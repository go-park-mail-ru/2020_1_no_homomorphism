#!/bin/bash
PROC="$(ps auxt| grep -i exe/server  | awk '{ print $2 }')"
kill -9  $PROC
PROC="$(ps auxt| grep -i exe/fileserver  | awk '{ print $2 }')"
kill -9  $PROC
