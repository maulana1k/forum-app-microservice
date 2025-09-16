import logging
import time
from typing import Any, Callable, Dict
import grpc
from grpc import ServicerContext
from google.protobuf.json_format import MessageToDict

logger = logging.getLogger(__name__)

class LoggingInterceptor(grpc.aio.ServerInterceptor):
    async def intercept_service(
        self,
        continuation: Callable,
        handler_call_details: grpc.HandlerCallDetails,
    ) -> grpc.RpcMethodHandler:
        # Log incoming call
        method = handler_call_details.method
        logger.info(f"gRPC call started: {method}")
        
        start_time = time.time()
        
        try:
            # Continue with the call
            response = await continuation(handler_call_details)
            return response
        except Exception as e:
            # Log error
            duration = time.time() - start_time
            logger.error(
                f"gRPC call failed: {method}, "
                f"error: {str(e)}, "
                f"duration: {duration:.3f}s"
            )
            raise
        finally:
            # Log successful completion
            duration = time.time() - start_time
            logger.info(
                f"gRPC call completed: {method}, "
                f"duration: {duration:.3f}s"
            )

class DetailedLoggingInterceptor(grpc.aio.ServerInterceptor):
    async def intercept_service(
        self,
        continuation: Callable,
        handler_call_details: grpc.HandlerCallDetails,
    ) -> grpc.RpcMethodHandler:
        method = handler_call_details.method
        start_time = time.time()
        
        logger.info(f"📡 gRPC Request received: {method}")
        
        try:
            response = await continuation(handler_call_details)
            duration = time.time() - start_time
            
            logger.info(
                f"✅ gRPC Request succeeded: {method}, "
                f"duration: {duration:.3f}s"
            )
            
            return response
            
        except Exception as e:
            duration = time.time() - start_time
            logger.error(
                f"❌ gRPC Request failed: {method}, "
                f"error: {type(e).__name__}: {str(e)}, "
                f"duration: {duration:.3f}s"
            )
            raise
    