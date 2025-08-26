from fastapi import APIRouter, Depends
from sqlalchemy.orm import Session
from sqlalchemy import text
from src.config.database import get_db

api_router = APIRouter()

@api_router.get("/", tags=["hello"])
def read_root():
    return {"message": "Welcome to the FastAPI ML Server!"}


@api_router.get("/health", tags=["health"])
def health_check(db: Session = Depends(get_db)):
    """
    Health check endpoint to verify database connection.
    """  
    try:
        db.execute(text("SELECT 1"))  # Executes a simple query to check DB connectivity
        return {"status": "healthy", "database": "connected"}
    except Exception as e:
        return {"status": "unhealthy", "database": "disconnected", "error": str(e)}
