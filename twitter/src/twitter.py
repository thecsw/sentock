from vaderSentiment.vaderSentiment import SentimentIntensityAnalyzer
import tweepy
import os
import time
import datetime

def tweet_in_company_tweets(elt, arr):
    for date in arr:
        for tweets in date:
            if elt in tweets:
                return True

    return False

#get tokens from environment var

consumer_token = 'XWYNu0pXy3FujGTFs5Zz2sI3w'
consumer_secret = 'u4lOvuenKT5rM5fvadgbe9MmRwUTK6NIcf6ZaB1HCRBKYX3lkB'
access_token = '1245747355397349386-jnDu79yI1J9UEKcFTaG1xV5HcvdZPo'
access_token_secret = 'nK1S3NvBI0Ki5bvKDpeZci2y5A3rk0hsimOcgrS5YyGki'

auth = tweepy.OAuthHandler(consumer_token, consumer_secret)
auth.set_access_token(access_token, access_token_secret)

api = tweepy.API(auth)

#this is the thing we want to search for in the tweets e.g. stock name
query = 'McDonald\'s&-filter:retweets'
#number of tweets we want with this query in it
#we may just want while true and circle through queries and then have
# to wait it out becuase of rate limiting
num_of_tweets = 100
tweets = []
#for the purposes of not going over older tweets
past_tweets = []
#The end data structure we will be working with. 3D array tweet_data[i][j][k], with i being the company number, j being the day that is j days before today, and k being the index of individual tweets
tweet_data=[]

#To add a company, add the string name and then..
companies = ["McDonald\'s", "Chipotle", "Microsoft", "FedEx", "Disney"]  

#add the latitude, longitude, and radius around the HQ you want to search
list_of_coords = ["41.8871,-87.6298,100mi", "33.628342,-117.927933,100mi", "47.673988,-122.121513,100mi", "35.117500,-89.971107,100mi", "34.180840,-118.308968,100mi"] 

#Get the dates we will search for
base = datetime.datetime.today()
date_list = [str(base - datetime.timedelta(days=6-x)).split(" ")[0] for x in range(7)]
print(date_list)

#look at docs, dunno how to make this work with rate limitting
for i in range(len(companies)):
    print('Company: ' + companies[i])
    company_tweets = []
    for date in date_list:
        print('Date: ' + date)
        date_tweets = []
        query = companies[i] + '&-filter:retweets'
        while len(date_tweets) < num_of_tweets:
            tweets_left = num_of_tweets - len(tweets)
            try:
                newer_tweets = api.search(q=query, count=tweets_left, result_type='recent', geocode=list_of_coords[i], until=date)
                count = 0
                for tweet in newer_tweets:
                    if not tweet_in_company_tweets(tweet.text,company_tweets):
                        date_tweets.append(tweet.text)
                        count += 1
                        if len(tweets) >= num_of_tweets:
                            break
                company_tweets.append(date_tweets)
                if count == 0:
                    print("no tweets");
                    break
                else:
                    print("Found " + str(count) + " tweet(s)")
                    
            except tweepy.TweepError as e:
                time.sleep(900)

    tweet_data.append(company_tweets)

for company in company_tweets:
    for day in company:
        for tweet in day:
            print(tweet)

#give this a sentence and it will spit out a dictionary
#that has different fields. The best one should probably
#be the compound
# ex: ans = analyzer(<some sentence>)
#     print(str(ans['compound']))
analyzer = SentimentIntensityAnalyzer()

sentiment_analysis = []
for tweet in tweets:
    sentiment_on_tweet = analyzer.polarity_scores(tweet.text)
    sentiment_analysis.append(sentiment_on_tweet)
