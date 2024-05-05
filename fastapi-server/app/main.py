# main.py
"""
The main module contains the FastAPI application instance
and the configuration for the application.
It imports the sentiment_router from the sentiment module
and includes it in the application using the include_router method.
The main module also defines a root route
that returns a simple message when accessed.
"""

from fastapi import FastAPI
from app.sentiment.routes import sentiment_router


app = FastAPI()
app.include_router(sentiment_router, prefix="/api")

if __name__ == "__main__":
    import uvicorn

    uvicorn.run(app, host="0.0.0.0", port=8000)


@app.get("/")
async def root():
    """
    The root route returns a simple message when accessed."""
    return {"message": "Hello World"}
