use eframe::NativeOptions;
use obsidian_android_view_app::DemoApp;

fn main() -> Result<(), eframe::Error> {
    let options = NativeOptions::default();
    DemoApp::run(options)
}
