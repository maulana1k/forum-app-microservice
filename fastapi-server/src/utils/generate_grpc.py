#!/usr/bin/env python3

import os
import subprocess

def generate_grpc_code():
    proto_file = "proto/recommender.proto"
    output_dir = "src/grpc_service"
    
    # Create output directory if it doesn't exist
    os.makedirs(output_dir, exist_ok=True)
    
    # Generate gRPC code
    cmd = [
        "python", "-m", "grpc_tools.protoc",
        f"--proto_path=proto",
        f"--python_out={output_dir}",
        f"--grpc_python_out={output_dir}",
        f"--pyi_out={output_dir}",
        proto_file
    ]
    
    result = subprocess.run(cmd, capture_output=True, text=True)
    
    if result.returncode != 0:
        print(f"Error generating gRPC code: {result.stderr}")
    else:
        print("gRPC code generated successfully!")
        
        # Fix import issues in generated code
        fix_imports(output_dir)

def fix_imports(output_dir):
    """Fix import statements in generated files"""
    files_to_fix = ["recommender_pb2.py", "recommender_pb2_grpc.py"]
    
    for filename in files_to_fix:
        filepath = os.path.join(output_dir, filename)
        if os.path.exists(filepath):
            with open(filepath, 'r') as f:
                content = f.read()
            
            # Fix import statements
            content = content.replace(
                "import recommender_pb2 as recommender__pb2",
                "from src.grpc_service import recommender_pb2 as recommender__pb2"
            )
            
            with open(filepath, 'w') as f:
                f.write(content)

if __name__ == "__main__":
    generate_grpc_code()