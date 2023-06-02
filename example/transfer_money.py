import uuid
import requests
import asyncio
import aiohttp


class AException(Exception):
    pass


async def query_A(session, account_uuid, url):
    request_id = uuid.uuid4()
    print(f"Make request for account={account_uuid} with request_id={uuid.uuid4()}")
    resp = await session.post(url.format(account_uuid=account_uuid, request_id=request_id), data=str(1 + index * 10))
    if resp.status != 200:
        raise AException


async def query_A_20_times(account_uuid, url):
    async with aiohttp.ClientSession(raise_for_status=True) as session:
        tasks = [asyncio.create_task(query_A(session, account_uuid, url)) for _ in range(20)]
        done, _ = await asyncio.wait(tasks)


account_uuids = [
    'e85bb517-8bc3-4939-bbaa-a080ce68d6d2',
    '9905260b-b7ba-459a-b90c-62dde56209f5',
    '71ed2c34-b613-4ab0-b201-bc709f3471bc',
    '34782b37-76df-4478-a325-fe61bfd039fc',
    '1cf8fd83-36a4-4c44-a671-798bdd3b8936',
    'f80ba086-1df5-4884-9a57-a85b5bc1098f',
]
url = "http://localhost:8080/transfer-money/{account_uuid}/{request_id}"
session = requests.Session()

for index, account_uuid in enumerate(account_uuids[:-1]):
    request_id = uuid.uuid4()
    print(f"Make request for account={account_uuid} with request_id={uuid.uuid4()}")
    resp = session.post(url.format(account_uuid=account_uuid, request_id=request_id), data=str(1 + index * 10))
    resp.raise_for_status()

# check 10 rps per user - should raise error
asyncio.run(query_A_20_times(account_uuids[-1], url))
