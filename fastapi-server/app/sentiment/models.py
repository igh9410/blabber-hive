# app/sentiment/models.py
"""
The models module contains Pydantic models that represent the request
and response bodies for the sentiment analysis endpoint.
The SentimentRequest model has a single field, text,
which is a string representing the text to be analyzed.
The SentimentResponse model has a single field, sentiment,
which is a string representing the sentiment of the text.
"""
from pydantic import BaseModel


class SentimentRequest(BaseModel):
    """
    The SentimentRequest model is a Pydantic model
    that represents the request body for the sentiment analysis endpoint.
    It has a single field, text, which is a string representing
    the text to be analyzed.
    he text field is required, meaning
    that it must be present in the request body.
    """

    text: str


class SentimentResponse(BaseModel):
    """
    The SentimentResponse model is a Pydantic model
    that represents the response body for the sentiment analysis endpoint.
    It has a single field, sentiment, which is
    a string representing the sentiment of the text.
    The sentiment field is required,
    meaning that it must be present in the response body.
    """

    sentiment: str
