import asyncio
from websockets.asyncio.server import serve
import time

async def send(websocket):
    await websocket.send('test message')
    time.sleep(5)

async def main():
    async with serve(send, 'localhost', 8765) as server:
        await server.serve_forever()

if __name__ == '__main__':
    asyncio.run(main())