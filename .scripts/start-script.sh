#!/bin/bash

# This script should kill all running instances (because init.d and CodeDeploy
# scripts may run multiple instances) and start a single new BTP process.
# This should ideally be run when the instance is updated with new BTP code.

# Store the pids of any running BTP instances
btpPids=$(pidof /home/ec2-user/app/bott-the-pigeon)

# Iterate through pids if it's not empty (so we're not trying to kill null)
if ((${btpPids[@]})); then
        for each in "${btpPids[@]}"
        do
                # Kill process using SIGTERM
                kill -15 $each
        done
fi

# Run a single BTP process in the background, using the production flag.
/home/ec2-user/app/bott-the-pigeon --prod > /dev/null 2> /dev/null < /dev/null &