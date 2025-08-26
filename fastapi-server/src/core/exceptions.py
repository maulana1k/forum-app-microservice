# src/core/exceptions.py
from fastapi import FastAPI, Request
from fastapi.responses import JSONResponse

class CustomAPIException(Exception):
    def __init__(self, message: str, status_code: int = 400):
        self.message = message
        self.status_code = status_code

def setup_exception_handlers(app: FastAPI):
    @app.exception_handler(CustomAPIException)
    async def custom_api_exception_handler(request: Request, exc: CustomAPIException):
        return JSONResponse(
            status_code=exc.status_code,
            content={"error": exc.message},
        )

    @app.exception_handler(Exception)
    async def global_exception_handler(request: Request, exc: Exception):
        return JSONResponse(
            status_code=500,
            content={"error": "Internal Server Error"},
        )
