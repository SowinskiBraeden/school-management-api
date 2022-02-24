#!/usr/bin/python3
import random
import names
from courses import courses

mockStudents = []

# Generate n students for mock data
def generateMockStudents(n):
  for _ in range(n):
    newStudent = {
      "name": names.get_full_name(),
      "requests": [], # list of class codes
      "schedule": {
        "block1": "",
        "block2": "",
        "block3": "",
        "block4": "",
        "block5": "",
        "block6": "",
        "block7": "",
        "block8": ""
      }
    }
    # Get list of random class choices with no repeats
    # 8 primary choices, 2 secondary choices
    courseSelection = random.sample(range(0, len(courses)), 10)
    for courseNum in courseSelection:
      newStudent["requests"].append(list(courses)[courseNum])
    mockStudents.append(newStudent)

  return mockStudents
