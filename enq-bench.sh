#!/usr/bin/env bash

# Prepare two branches <old> and <new>, both with your benchmark test yaml located
# at scheduler_perf/config/performance-config.yaml
# <new> branch should contain the EventsToRegister() implementation for your plugin.

for br in churn-cluster-op churn-cluster-op-zsm; do
  git checkout $br
  echo "Switched to $br"
  # Run the test 10 times.
  for i in {1..10}; do
    # <1000 5000> may need to be tweaked depending on your benchmark test name.
    for nodes in 1000 5000; do
      echo "====$i-${node}Nodes===="
      # Note: you need to replace SchedulingWithPodChurn/${nodes}Nodes with your
      # corresponding (sub)benchmark test name.
      make test-integration WHAT=./test/integration/scheduler_perf KUBE_TEST_VMODULE="''" KUBE_TEST_ARGS="-alsologtostderr=false -logtostderr=false -run=^$$ -benchtime=1ns -bench=BenchmarkPerfScheduling/SchedulingWithNodeLabelChurn/${nodes}Nodes -data-items-dir ~/logs/$br/$nodes"
      sleep 5
    done
  done
  echo
done

# After the above completes, it's expected to get a bunch of Benchmark*.json
# located at ~/logs/{old|new}/{1000|5000}.
# Next, let's concat them into a single file
for folder in ~/logs/churn-cluster-op/1000 ~/logs/churn-cluster-op/5000 ~/logs/churn-cluster-op-zsm/1000 ~/logs/churn-cluster-op-zsm/5000; do
  ls $folder/BenchmarkPerfScheduling* | while read f; do
    echo "===$f==="
    cat $f >> $folder/result.txt
    echo "" >> $folder/result.txt
  done
done

# If you see a lot of timeout issues and hence cannot get desired results (a lot of
# BenchmarkPerfScheduling*.json gets a single line with null result).
# Try to rebase #96696 to work it around.

# Next, you should be able to compare ~/logs/old/{1000|5000}/result.txt with
# ~/logs/new/{1000|5000}/result.txt
# You can leverage https://github.com/Huang-Wei/k8s-sched-perf-stat to get the diff in one command.