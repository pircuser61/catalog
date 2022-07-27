protoc ^
-I . ^
--go_out=. ^
--go_opt=paths=source_relative ^
--go-grpc_out=. ^
--go-grpc_opt=paths=source_relative ^
--grpc-gateway_out . ^
--grpc-gateway_opt paths=source_relative ^
--grpc-gateway_opt allow_delete_body=true ^
--openapiv2_out . ^
--openapiv2_opt allow_delete_body=true ^
.\api\api.proto 



