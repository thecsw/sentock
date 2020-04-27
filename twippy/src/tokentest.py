#!/bin/python
import tweepy
import os

# open_tweepy attempts to open a tweepy connection
def open_tweepy(consumer_token, consumer_secret, access_token, access_token_secret):
    try:
        auth = tweepy.OAuthHandler(consumer_token, consumer_secret)
        auth.set_access_token(access_token, access_token_secret)
        api = tweepy.API(auth)

        print(api.me())
        return 1
    except Exception as e:
        print(e)
        return 0


def test_tweepy():
    # Invalid tweet
    assert open_tweepy("INVALID", "INVALID", "INVALID", "INVALID") == 0
    # Valid tweet
    consumer_token = os.environ["CONSUMER_TOKEN"]
    consumer_secret = os.environ["CONSUMER_SECRET"]
    access_token = os.environ["ACCESS_TOKEN"]
    access_token_secret = os.environ["ACCESS_TOKEN_SECRET"]
    assert (
        open_tweepy(consumer_token, consumer_secret, access_token, access_token_secret)
        == 1
    )


if __name__ == "__main__":
    test_tweepy()
