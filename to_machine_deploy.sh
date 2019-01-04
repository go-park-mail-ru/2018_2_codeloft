#!/bin/bash

export SSHPASS=$DEPLOY_SSH
sshpass -e ssh -o StrictHostKeychecking=no $DEPLOY_USER@$DEPLOY_HOST "cd $DEPLOY_PATH && ./deploy.sh $TRAVIS_BRANCH"