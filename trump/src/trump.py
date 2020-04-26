#!/bin/python
import tweepy
import json
import time
import requests
import json
import os
from vaderSentiment.vaderSentiment import SentimentIntensityAnalyzer

#test out tweet streaming:
consumer_token = 'XWYNu0pXy3FujGTFs5Zz2sI3w'
consumer_secret = 'u4lOvuenKT5rM5fvadgbe9MmRwUTK6NIcf6ZaB1HCRBKYX3lkB'
access_token = '1245747355397349386-jnDu79yI1J9UEKcFTaG1xV5HcvdZPo'
access_token_secret = 'nK1S3NvBI0Ki5bvKDpeZci2y5A3rk0hsimOcgrS5YyGki'

auth = tweepy.OAuthHandler(consumer_token, consumer_secret)
auth.set_access_token(access_token, access_token_secret)

api = tweepy.API(auth)

#sentiment analyzer stuff:
analyzer = SentimentIntensityAnalyzer()

texts = []

URL = 'http://server:10000/sentiments'
session = requests.session()

# Some aliases came from
# https://en.wikipedia.org/wiki/List_of_nicknames_of_presidents_of_the_United_States#Donald_Trump
# https://en.wikipedia.org/wiki/List_of_nicknames_used_by_Donald_Trump
# http://www.citypages.com/news/the-20-best-nicknames-for-donald-trump-so-far/389377462
# https://www.cheatsheet.com/entertainment/15-hilarious-nicknames-donald-trump-has-been-called.html/
# https://medium.com/@allanishac/50-nicknames-for-donald-trump-you-wont-be-hearing-on-fox-news-7b7a5ca1b1b1
companies = ["Trump",
             "Trump\'s",
             "POTUS",
             "White House President",
             "President",
             "Donald Trump",
             "Commander-in-chief",
             "President of the United States",
             "President of United States",
             "US President",
             "President of US",
             "Leader of the free world",
             "The Donald",
             "45, the 45",
             "Conspiracy Theorist-in-Chief",
             "President Snowflake",
             "Snowflake-in-Chief",
             "Donald John Trump",
             "President T.",
             "Your Favorite President",
             "Prima Donald",
]

def got_tweet(tweet_id, text, created_at):
    sentimentValue = float(analyzer.polarity_scores(text)["compound"])
    if (sentimentValue == 0):
        return "No sentiment value"
    payload = {
        "company":   "Trump",
        "tweet_id":  tweet_id,
        "sentiment": float(sentimentValue),
        "unix":      int(created_at),
    }
    headers = {
        "webkey": os.environ["WEB_KEY"]
    }
    # Make the call
    requests.post(URL, data=json.dumps(payload), headers=headers)

class SentStreamListener(tweepy.StreamListener):

    def on_status(self, status):
        # Get the actual payload
        msg = status._json
        # make sure it is not a retweet (we ignore retweets)
        if 'retweeted_status' in msg:
            return True
        # Log the tweet's creation time in unix timestamp
        createdAt = int(time.mktime(time.strptime(msg["created_at"].replace("+0000",""))))
        # If the tweet text is collapsed because its too long:
        if msg['truncated']:
            got_tweet(msg['id_str'], msg['extended_tweet']['full_text'], createdAt)
            return True
        # The tweet is a regular one (not a retweet or more than 140? character long)
        got_tweet(msg['id_str'], msg['text'], createdAt)
        return True
 
    # Callback on error
    def on_error(self, status_code):
        if status_code == 420:
            #returning False in on_error disconnects the stream
            print("Error code 420 occurred, waiting 15 seconds before continuing", flush=True)
            time.sleep(15)
            return True
        print("Uncaught Error occured while listening: " + str(status_code), flush=True)
        return False

if __name__=="__main__":
    print("POTUS THE 45 SENTIMENT ANALYSIS")
    print("Gratiously waiting until other microservices become operational...")
    time.sleep(5)
    print("Setting up listener...")
    sentStreamListener = SentStreamListener()
    mstream = tweepy.Stream(auth = api.auth, listener=sentStreamListener)
    mstream.filter(track=companies, is_async=False)
    print("waiting...")
    while True:
        time.sleep(1800)
    print("Disconnecting..")
    mstream.disconnect()
