package main

import (
	"context"
	"fmt"
	"image/color"
	"log"
	"os"
	"time"

	"gioui.org/app"
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/text"
	"gioui.org/unit"
	"gioui.org/widget"
	"gioui.org/widget/material"
	"github.com/coder/websocket"
)

type ObsidianApp struct {
	window      *app.Window
	urlInput    widget.Editor
	url         string
	startButton widget.Clickable
	conn        *websocket.Conn
	wsCtx       context.Context
}

type C = layout.Context
type D = layout.Dimensions

func main() {
	obApp := ObsidianApp{}
	obApp.url = ""

	go func(oapp *ObsidianApp) {
		oapp.window = new(app.Window)
		err := oapp.Run(oapp.window)

		if err != nil {
			log.Fatal(err)
		}

		os.Exit(0)
	}(&obApp)

	app.Main()
}

func (o *ObsidianApp) Run(window *app.Window) error {
	theme := material.NewTheme()

	var ops op.Ops

	for {
		switch e := window.Event().(type) {
		case app.DestroyEvent:
			return e.Err
		case app.FrameEvent:
			// This graphics context is used for managing the rendering state.
			gtx := app.NewContext(&ops, e)

			if o.url == "" {
				if o.startButton.Clicked(gtx) {
					o.url = o.urlInput.Text()
				} else {
					layout.Flex{
						Axis:    layout.Vertical,
						Spacing: layout.SpaceStart,
					}.Layout(
						gtx,
						// The title
						layout.Rigid(
							func(gtx C) D {
								title := material.H2(theme, "Obsidian Android View")
								maroon := color.NRGBA{R: 127, G: 0, B: 0, A: 255}
								title.Color = maroon
								title.Alignment = text.Middle

								return title.Layout(gtx)
							},
						),
						// Spacer
						layout.Rigid(
							layout.Spacer{Height: unit.Dp(400)}.Layout,
						),
						// The input
						layout.Rigid(
							func(gtx C) D {
								ed := material.Editor(theme, &o.urlInput, "IP Address")
								o.urlInput.SingleLine = true
								o.urlInput.Alignment = text.Middle
								return ed.Layout(gtx)
							},
						),
						// The button
						layout.Rigid(
							func(gtx C) D {
								margins := layout.Inset{
									Top:    unit.Dp(25),
									Bottom: unit.Dp(25),
									Right:  unit.Dp(35),
									Left:   unit.Dp(35),
								}

								return margins.Layout(gtx,
									func(gtx C) D {
										btn := material.Button(theme, &o.startButton, "Start")
										return btn.Layout(gtx)
									},
								)
							},
						),
					)
				}
			} else {
				if o.conn == nil {
					ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
					defer cancel()
					o.wsCtx = ctx
					c, _, err := websocket.Dial(ctx, "ws://"+o.url+":8080", &websocket.DialOptions{
						Subprotocols: []string{"send"},
					})

					if err != nil {
						log.Fatal(err)
					}

					o.conn = c
					defer o.conn.Close(websocket.StatusInternalError, "the sky is falling")
				} else {
					_, message, err := o.conn.Read(o.wsCtx)

					if err != nil {
						log.Fatal(err)
					}

					fmt.Printf("recv: %s\n", message)
				}
			}

			// Pass the drawing operations to the GPU.
			e.Frame(gtx.Ops)
		}
	}
}
