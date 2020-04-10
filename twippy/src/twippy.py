#!/bin/python
import tweepy
import json
import time
import databa
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

# List of companies we would like to find in our tweets
companies = ["McDonalds", "McDonald\'s", "Chipotle", "Chipotle\'s", "Microsoft", "Microsoft\'s", "Fedex", "Fedex\'s", "Disney", "Walt Disney", "Disney\'s"]

# Create the main table
databa.create_table()

def got_tweet(tweet_id, text, created_at):
    sentimentValue = float(analyzer.polarity_scores(text)["compound"])
    if (sentimentValue == 0):
        return
    # Compute the company
    company = ""
    if companies[0].lower() in text.lower() or companies[1].lower() in text.lower():
        company = "McDonalds"
    elif companies[2].lower() in text.lower() or companies[3].lower() in text.lower():
        company = "Chipotle"
    elif companies[4].lower() in text.lower() or companies[5].lower() in text.lower():
        company = "Microsoft"
    elif companies[6].lower() in text.lower() or companies[7].lower() in text.lower():
        company = "Fedex"
    elif companies[8].lower() in text.lower() or companies[9].lower() in text.lower() or companies[10].lower() in text.lower():
        company = "Disney"
    if company == "":
        return
    print(f"{text} | sentiment: {sentimentValue}", flush=True)
    databa.add_sentiment(
        company, 
        tweet_id, 
        text, 
        created_at, 
        sentimentValue)

class TestStreamListener(tweepy.StreamListener):

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

print("Setting up listener...")
testStreamListener = TestStreamListener()
mstream = tweepy.Stream(auth = api.auth, listener=testStreamListener)
mstream.filter(track=companies, is_async=True)
print("waiting...")
while True:
    time.sleep(1800)
print("Disconnecting..")
mstream.disconnect()
databa.close()