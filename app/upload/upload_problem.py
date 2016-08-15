"""
[problem]
title=xxx
description=xxx
time=1
memory=100
input=in
output=out
simpleinput=sinput
simpleoutput=soutput
solved=0
display=true/false
"""
#encoding:utf-8

import sys
import ConfigParser

if len(sys.argv) != 2:
	print 'Usage: python upload_problem.py problem.in'
	sys.exit()
filename = sys.argv[1]

config = ConfigParser.ConfigParser()
config.read(filename)
title = config.get('problem', 'title')
description = config.get('problem', 'description')
time = config.get('problem', 'time')
memory = config.get('problem', 'memory')
_input = config.get('problem', 'input')
output = config.get('problem', 'output')
simpleinput = config.get('problem', 'simpleinput')
simpleoutput = config.get('problem', 'simpleoutput')
solved = 0
display = False
