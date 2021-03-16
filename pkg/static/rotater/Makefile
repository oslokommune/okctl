package:
	docker build . --iidfile image-id
	docker run --rm -d -v $$(pwd):/host $$(cat image-id) cp /source/lambda_function.zip /host/.
	echo $!