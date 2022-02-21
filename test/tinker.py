#!/usr/bin/python3
import random
import names
import math
import json
from courses import courses, activeCourses

'''
  I will be using python to test and
  develop course selection / schedule
  generation

  this is honestly terrible code but 
  this is not an easy task for me to
  do
'''

'''
schedule example:
schedule: {
  "semester1": {
    "block1": "className",
    "block2": "className",
    "block3": "className",
    "block4": "className" 
  },
  "semester2": {
    "block1": "className",
    "block2": "className",
    "block3": "className",
    "block4": "className"
  }
}
'''

'''
running example:
running: {
  "semester1": {
    "block1": {classCode:{"className":name,"students":[student Name]},classCode:{"className":name,"students":[student Name]}},
    "block2": {classCode:{"className":name,"students":[student Name]},classCode:{"className":name,"students":[student Name]}},
    "block3": {classCode:{"className":name,"students":[student Name]},classCode:{"className":name,"students":[student Name]}},
    "block4": {classCode:{"className":name,"students":[student Name]},classCode:{"className":name,"students":[student Name]}}
  },
  "semester2": {
    "block1": {classCode:{"className":name,"students":[student Name]},classCode:{"className":name,"students":[student Name]}},
    "block2": {classCode:{"className":name,"students":[student Name]},classCode:{"className":name,"students":[student Name]}},
    "block3": {classCode:{"className":name,"students":[student Name]},classCode:{"className":name,"students":[student Name]}},
    "block4": {classCode:{"className":name,"students":[student Name]},classCode:{"className":name,"students":[student Name]}}
  ]
}
'''

minReq, classCap, blockClassLimit = 18, 30, 15
mockStudents = []
running = {
  "semester1": {
    "block1": {},
    "block2": {},
    "block3": {},
    "block4": {}
  },
  "semester2": {
    "block1": {},
    "block2": {},
    "block3": {},
    "block4": {}
  }
}


# Generate n students for mock data
def generateMockStudents(n):
  for _ in range(n):
    newStudent = {
      "name": names.get_full_name(),
      "requests": [], # list of class codes
      "schedule": {
        "semester1": {
          "block1": "",
          "block2": "",
          "block3": "",
          "block4": ""
        },
        "semester2": {
          "block1": "",
          "block2": "",
          "block3": "",
          "block4": ""
        }
      }
    }
    # Get list of random class choices with no repeats
    # 8 primary choices, 2 secondary choices
    courseSelection = random.sample(range(0, len(courses)), 10)
    for courseNum in courseSelection:
      newStudent["requests"].append(list(courses)[courseNum])
    mockStudents.append(newStudent)


def generateSchedule():
  # Collect data and calculate schedules
  for student in mockStudents:
    # Tally class request
    for request in student["requests"]:
      courses[request]["totalrequests"] += 1
      courses[request]["studentindexes"].append(mockStudents.index(student))
      # Add course to active list if enough requests      
      if courses[request]["totalrequests"] > minReq and courses[request]["code"] not in activeCourses: activeCourses[courses[request]["code"]] = courses[request]

  # calculate # of times to run class
  for i in range(len(activeCourses)):
    classRunCount = math.floor(activeCourses[list(activeCourses)[i]]["totalrequests"] / classCap)
    # If there is 18+ requests left, 1 more class could be run
    if (activeCourses[list(activeCourses)[i]]["totalrequests"] % classCap) > minReq: classRunCount += 1
    activeCourses[list(activeCourses)[i]]["classRunCount"] = classRunCount

  for student in mockStudents:
    alternateOffset = len(student["requests"])-8
    for i in range(len(student["requests"])-alternateOffset): # Subtract x classes as they are alternatives
      currentCourse = student["requests"][i]
      while True:
        if currentCourse in activeCourses:
          blockIndex = 1
          semesterIndex = 1
          while True:
            block = f"block{blockIndex}"
            semester = f"semester{semesterIndex}"
            if student["schedule"][semester][block] != "":
              if semesterIndex == 2 and blockIndex == 4:
                # No available classes
                print("No more available classes")
                break
              elif blockIndex == 4:
                blockIndex = 1
                semesterIndex += 1
              else: blockIndex += 1
            if currentCourse in running[semester][block]:
              # Add student to class
              if len(running[semester][block][currentCourse]["students"]) < classCap:
                running[semester][block][currentCourse]["students"].append(student["name"])
                student["schedule"][semester][block] = running[semester][block][currentCourse]["name"]
                break
              else: # Find next available class or create new one
                if semesterIndex == 2 and blockIndex == 4:
                  # No available classes
                  print("No more available classes")
                  break
                elif blockIndex == 4:
                  blockIndex = 1
                  semesterIndex += 1
                else: blockIndex += 1
            else:
              if semesterIndex == 2 and blockIndex == 4:
                # Class does not exists
                # Create new class in first available slot 
                blockNum = 1
                semesterNum = 1
                while True:
                  newBlock = f"block{blockNum}"
                  newSemester = f"semester{semesterNum}"
                  if currentCourse not in running[newSemester][newBlock] and len(running[newSemester][newBlock]) < blockClassLimit:
                    running[newSemester][newBlock][currentCourse] = {
                      "name": courses[currentCourse]["name"],
                      "students": [student["name"]]
                    }
                    break
                  else:
                    if semesterNum == 2 and blockNum == 4:
                      # No room in school for more classes
                      print("No room in school for more classes")
                      break
                    elif blockNum == 4:
                      blockNum = 1
                      semesterNum += 1
                    else: blockNum += 1
                break
              elif blockIndex == 4: # 4th Block (final block)
                blockIndex = 1
                semesterIndex += 1
              else: blockIndex += 1
          break
        elif currentCourse not in activeCourses:
          if alternateOffset == 0:
            # TODO: No more alternatives, solve problem
            print("No more alternatives, solve this problem somehow")
            break
          currentCourse = student["requests"][len(student["requests"])-alternateOffset]
          alternateOffset -= 1


if __name__ == '__main__':
  generateMockStudents(400)
  generateSchedule()

  with open("schedule.json", "w") as outfile:
    json.dump(running, outfile, indent=2)
