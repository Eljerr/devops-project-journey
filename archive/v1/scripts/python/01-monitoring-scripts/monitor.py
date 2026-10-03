import json
import logging
import argparse
import asyncio
import aiohttp


parser = argparse.ArgumentParser()
parser.add_argument("--output", required=True)
args = parser.parse_args()

logging.basicConfig(
    level=logging.INFO, format="%(asctime)s - %(levelname)s - %(message)s"
)

target_service = [
    {"name": "Nginx Web", "url": "https://httpstat.us/200"},
    {"name": "NodeJS API", "url": "https://httpstat.us/502"},
    {"name": "Payment Gateway", "url": "https://domain-ini-pasti-error-123.com"},
]


class ServiceMonitor:
    def __init__(self, services, output_file):
        self.services = services
        self.output_file = output_file

    async def fetch_status(self, session, target):
        try:
            async with session.get(target["url"]) as response:
                if response.status == 200:
                    return {"name": target["name"], "status": "UP"}
                else:
                    logging.warning(
                        f"Perlu tindak lebih lanjut pada server {target['name']}"
                    )
                    return {"name": target["name"], "status": "DOWN"}
        except aiohttp.ClientError as e:
            logging.error(f"Server {target['name']} tidak merespon")
            return {"name": target["name"], "status": "ERROR", "detail": str(e)}

    async def run(self):
        logging.info("Memulai pengecekan kesehatan server secara asinkron...")

        async with aiohttp.ClientSession() as session:
            tasks = [self.fetch_status(session, target) for target in self.services]

            result = await asyncio.gather(*tasks)

            with open(self.output_file, "w") as file:
                json.dump(result, file, indent=4)


monitor = ServiceMonitor(target_service, args.output)
asyncio.run(monitor.run())
