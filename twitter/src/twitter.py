from vaderSentiment.vaderSentiment import SentimentIntensityAnalyzer
import tweepy
import os
import time

#get tokens from environment var

consumer_token = os.environ['CONSUMER_TOKEN']
consumer_secret = os.environ['CONSUMER_SECRET']
access_token = os.environ['ACCESS_TOKEN']
access_token_secret = os.environ['ACCESS_SECRET']

auth = tweepy.OAuthHandler(consumer_token, consumer_secret)
auth.set_access_token(access_token, access_token_secret)

api = tweepy.API(auth)

#this is the thing we want to search for in the tweets e.g. stock name
query = ''
#number of tweets we want with this query in it
#we may just want while true and circle through queries and then have
# to wait it out becuase of rate limiting
num_of_tweets = 100
tweets = []


#look at docs, dunno how to make this work with rate limitting
while len(tweets) < num_of_tweets:
    tweets_left = num_of_tweets - len(tweets)
    try:
        newer_tweets = api.search(q=query)
        if len(newer_tweets) == 0:
            break
        tweets.extend(newer_tweets)
    except tweepy.TweepError as e:
        #figure out how long to sleep for
        time.sleep(120)


#give this a sentence and it will spit out a dictionary
#that has different fields. The best one should probably
#be the compound
# ex: ans = analyzer(<some sentence>)
#     print(str(ans['compound']))
analyzer = SentimentIntensityAnalyzer()

sentiment_analysis = []
for tweet in tweets:
    sentiment_on_tweet = analyzer(tweet)
    sentiment_analysis.append(sentiment_on_tweet)
