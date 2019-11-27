docker build $docker_build_directory -t $docker_repo:$tmpname -f $dockerfile

for tag in ${tags//,/ }
do
  docker tag $docker_repo:$tmpname $docker_repo:$tag
  docker push $docker_repo:$tag
done
