#!/usr/bin/env bash

godepgraph -s github.com/bgrewell/dtac-agent/... > /tmp/deps.dot
grep -E 'github.com/intel-innersource/|digraph|}' /tmp/deps.dot > /tmp/filtered_deps.dot
dot -Tpng /tmp/filtered_deps.dot -o /tmp/graph.png
xdg-open /tmp/graph.png