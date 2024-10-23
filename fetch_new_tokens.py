import time

import requests
from loguru import logger

processed_event_ids = set()


def fetch_events():
    try:
        response = requests.get("https://api.trongrid.io/v1/contracts/TTfvyrAz86hbZk5iDpKD78pqLGgi8C7AAw/events")
        response.raise_for_status()
        return response.json()
    except requests.RequestException as e:
        logger.info(f"Error: {e}")


def get_event_by_index(transaction_id, index):
    url = f"https://api.trongrid.io/v1/transactions/{transaction_id}/events"
    response = requests.get(url)
    if response.status_code == 200:
        return next((e for e in response.json().get('data', []) if e.get('event_index') == index), None)
    logger.info(f"Error fetching data: {response.status_code}")


def check_events_for_token_create(events):
    if not events or 'data' not in events:
        logger.info("No data")
        return

    for event in events['data']:
        if event.get('event_name') == 'TokenCreate' and event.get('transaction_id') not in processed_event_ids:
            logger.info(f"Event TokenCreate: {event}")
            token = get_event_by_index(event['transaction_id'], 0)
            if token:
                logger.info(f"Token: {token['caller_contract_address']}")
            processed_event_ids.add(event['transaction_id'])


def main():
    while True:
        check_events_for_token_create(fetch_events())
        time.sleep(1)


if __name__ == "__main__":
    main()
