import tweepy
import json
import time
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

class TestStreamListener(tweepy.StreamListener):

    def on_status(self, status):
        createdAt = time.strptime(status.created_at.replace("+0000", ""))
        print("New tweet found, extended: " + str(status.truncated) + 
              ", Retweet: " + str(hasattr(status, 'retweeted_status'))+
              ", Created At " + str(createdAt))
        # The tweet is a retweet
        if hasattr(status, 'retweeted_status'):
            if status.retweeted_status['truncated']:
                texts.append((status.retweeted_status['extended_tweet']['full_text'], createdAt))
            return True
        # The tweet text is collapsed because of its length
        if status.truncated:
            #print(status.extended_tweet['full_text'])
            texts.append((status.extended_tweet['full_text'], createdAt))
            return True
        # The tweet is a regular one (not a retweet or more than 140? character long)
        print(status.text)
        texts.append((status.text, createdAt))
        return True

    # Callback on error
    def on_error(self, status_code):
        print("Error occured while listening: " + status_code)
        if status_code == 420:
            #returning False in on_error disconnects the stream
            return False
        return False

print("Setting up listener...")
testStreamListener = TestStreamListener()
mstream = tweepy.Stream(auth = api.auth, listener=testStreamListener)
mstream.filter(track=['McDonald,McDonalds,McDonald\'s'], is_async=True)
print("waiting...")
time.sleep(30)
print("Disconnecting..")
mstream.disconnect()
print("Number of texts gathered: " + str(len(texts)))
print("Analyzing sentiment...")
for text in texts:
    print(str(analyzer.polarity_scores(text[0])) + " : " + text[0])