import asyncio
from websockets.asyncio.server import serve
from websockets import broadcast

connected_clients = set()

async def handler(websocket):
    connected_clients.add(websocket)

    try:
        async for _ in websocket:
            pass # Keep the connection open
    finally:
        connected_clients.remove(websocket)

async def send_messages():
    while True:
        await asyncio.sleep(5)
        broadcast(connected_clients, 'test message')
        print('Sent')

async def main():
    server = await serve(handler, 'localhost', 8765)
    await asyncio.gather(server.wait_closed(), send_messages())

if __name__ == '__main__':
    asyncio.run(main())