#!/bin/python
import tweepy
import json
import time
import requests
import json
import os
from vaderSentiment.vaderSentiment import SentimentIntensityAnalyzer

# test out tweet streaming:
CONSUMER_TOKEN = "XWYNu0pXy3FujGTFs5Zz2sI3w"
CONSUMER_SECRET = "u4lOvuenKT5rM5fvadgbe9MmRwUTK6NIcf6ZaB1HCRBKYX3lkB"
ACCESS_TOKEN = "1245747355397349386-jnDu79yI1J9UEKcFTaG1xV5HcvdZPo"
ACCESS_TOKEN_SECRET = "nK1S3NvBI0Ki5bvKDpeZci2y5A3rk0hsimOcgrS5YyGki"

AUTH = tweepy.OAuthHandler(CONSUMER_TOKEN, CONSUMER_SECRET)
AUTH.set_access_token(ACCESS_TOKEN, ACCESS_TOKEN_SECRET)

API = tweepy.API(AUTH)

# sentiment analyzer stuff:
ANALYZER = SentimentIntensityAnalyzer()

TEXTS = []

URL = "http://server:10000/sentiments"
SESSION = requests.session()

KEYWORDS = {
    "McDonalds": "McDonalds",
    "McDonald's": "McDonalds",
    "Chipotle": "Chipotle",
    "Chipotle's": "Chipotle",
    "Microsoft": "Microsoft",
    "Microsoft's": "Microsoft",
    "Fedex": "Fedex",
    "Fedex's": "Fedex",
    "Disney": "Disney",
    "Disney's": "Disney",
    "Tesla": "Tesla",
    "TSLA": "Tesla",
    "Tesla's": "Tesla",
    "Google": "Google",
    '"Alphabet Inc"': "Google",
    "Google's": "Google",
    "Facebook": "Facebook",
    "Facebook's": "Facebook",
    "Amazon": "Amazon",
    "Amazon's": "Amazon",
}

# Get all unique companies for keywords
COMPANIES = list(set(KEYWORDS.values()))


def got_tweet(tweet_id, text, created_at):
    sentimentValue = float(ANALYZER.polarity_scores(text)["compound"])
    if sentimentValue == 0:
        return "No sentiment value"
    # Compute the company
    company = ""
    # Go through the keywords, where
    # the items() returns tuples, [0] is the keyword
    # and [1] is the company name
    for t in KEYWORDS.items():
        if t[0].lower() in text.lower():
            company = t[1]
            break
    if company == "":
        return "No company"
    print(f"{text} | sentiment: {sentimentValue}", flush=True)
    payload = {
        "company": company,
        "tweet_id": tweet_id,
        "sentiment": float(sentimentValue),
        "unix": int(created_at),
    }
    headers = {"webkey": os.environ["WEB_KEY"]}
    # Make the call
    requests.post(URL, data=json.dumps(payload), headers=headers)
    return company


class SentStreamListener(tweepy.StreamListener):
    def on_status(self, status):
        # Get the actual payload
        msg = status._json
        # make sure it is not a retweet (we ignore retweets)
        if "retweeted_status" in msg:
            return True
        # Log the tweet's creation time in unix timestamp
        createdAt = int(
            time.mktime(time.strptime(msg["created_at"].replace("+0000", "")))
        )
        # If the tweet text is collapsed because its too long:
        if msg["truncated"]:
            got_tweet(msg["id_str"], msg["extended_tweet"]["full_text"], createdAt)
            return True
        # The tweet is a regular one (not a retweet or more than 140? character long)
        got_tweet(msg["id_str"], msg["text"], createdAt)
        return True

    # Callback on error
    def on_error(self, status_code):
        if status_code == 420:
            # returning False in on_error disconnects the stream
            print(
                "Error code 420 occurred, waiting 5 seconds before continuing",
                flush=True,
            )
            time.sleep(5)
            return True
        print("Uncaught Error occured while listening: " + str(status_code), flush=True)
        return False


if __name__ == "__main__":
    print("Gratiously waiting until other microservices become operational...")
    time.sleep(5)
    print("Setting up listener...")
    sentStreamListener = SentStreamListener()
    mstream = tweepy.Stream(auth=API.auth, listener=sentStreamListener)
    mstream.filter(track=COMPANIES, is_async=False)
    print("waiting...")
    while True:
        time.sleep(1800)
    print("Disconnecting..")
    mstream.disconnect()
