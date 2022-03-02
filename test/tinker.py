#!/usr/bin/python3
import sys
import json
import math
from courses import courses, activeCourses
from mock import generateMockStudents

'''
  I will be using python to test and
  develop course selection / schedule
  generation

  this is honestly terrible code but 
  this is not an easy task for me to
  do

  Block 1-4 is first semester while
  block 5-8 is second semester
'''

'''
schedule example:
schedule: {
  "block1": "className",
  "block2": "className",
  "block3": "className",
  "block4": "className" 
  "block5": "className",
  "block6": "className",
  "block7": "className",
  "block8": "className"
}
'''

'''
running example:
running: {
  "block1": {classCode:{"className":name,"students":[student Name]},classCode:{"className":name,"students":[student Name]}},
  "block2": {classCode:{"className":name,"students":[student Name]},classCode:{"className":name,"students":[student Name]}},
  "block3": {classCode:{"className":name,"students":[student Name]},classCode:{"className":name,"students":[student Name]}},
  "block4": {classCode:{"className":name,"students":[student Name]},classCode:{"className":name,"students":[student Name]}},
  "block5": {classCode:{"className":name,"students":[student Name]},classCode:{"className":name,"students":[student Name]}},
  "block6": {classCode:{"className":name,"students":[student Name]},classCode:{"className":name,"students":[student Name]}},
  "block7": {classCode:{"className":name,"students":[student Name]},classCode:{"className":name,"students":[student Name]}},
  "block8": {classCode:{"className":name,"students":[student Name]},classCode:{"className":name,"students":[student Name]}}
}
'''

# Error 1: No classes in schedule can fit this student
# Error 2: No more room in schedule for another class

# 400 by default
studentsNum = 500

err1, err2, = 0,0
minReq, classCap, blockClassLimit = 18, 30, 18
mockStudents = []
running = {
  "block1": {},
  "block2": {},
  "block3": {},
  "block4": {},
  "block5": {},
  "block6": {},
  "block7": {},
  "block8": {}
}


def generateScheduleV1():
  global err1, err2
  # Collect data and calculate schedules
  for student in mockStudents:
    # Tally class request
    for request in student["requests"]:
      courses[request]["totalrequests"] += 1
      courses[request]["studentindexes"].append(mockStudents.index(student))
      # Add course to active list if enough requests
      if courses[request]["totalrequests"] > minReq and courses[request]["code"] not in activeCourses: activeCourses[courses[request]["code"]] = courses[request]

  for student in mockStudents:
    alternateOffset = len(student["requests"])-8
    alternateIndex = 8
    for i in range(len(student["requests"])-alternateOffset): # Subtract x classes as they are alternatives
      currentCourse = student["requests"][i]
      generate = True
      while generate:
        # If class is allowed to run
        if currentCourse in activeCourses:
          blockIndex = 1
          courseNum = 1
          getFreeBlock = True
          while getFreeBlock:
            block = f"block{blockIndex}"
            cname = f"{courses[currentCourse]['name']}-{courseNum}"
            if cname in running[block]:
              if student["schedule"][block] == "": # Add student to class
                if len(running[block][cname]["students"]) < classCap:
                  running[block][cname]["students"].append(student["name"])
                  student["schedule"][block] = cname
                  getFreeBlock = False
                else: # Find next available class or create new one
                  if blockIndex == 8:  # No available classes
                    if len(student["requests"]) == 8:
                      # This student never had any alternatives
                      # How to solve?

                      # print(f"\nError 1: No more available classes for student {student['name']}")
                      generate = False
                    else:
                      if alternateIndex <= (len(student["requests"]) - 1):
                        currentCourse = student["requests"][alternateIndex]
                        alternateIndex += 1
                      else: # No more alterantive
                        # print(f"\nError 1: No more available classes for student {student['name']}")
                        err1 += 1
                        generate = False
                    getFreeBlock = False
                  else:
                    if (courses[currentCourse]["teachers"] - courseNum) > 0:
                      courseNum += 1
                    else:
                      blockIndex += 1
                      courseNum = 1
              else:
                if blockIndex == 8: # No available classes
                  if len(student["requests"]) == 8:
                    # This student never had any alternatives
                    # How to solve?

                    # print(f"\nError 1: No more available classes for student {student['name']}")
                    err1 += 1
                    generate = False
                  else:
                    if alternateIndex <= (len(student["requests"]) - 1):
                      currentCourse = student["requests"][alternateIndex]
                      alternateIndex += 1
                    else: # No more alternatives
                      # print(f"\nError 1: No more available classes for student {student['name']}")
                      err1 += 1
                      generate = False
                  getFreeBlock = False
                else:
                  blockIndex += 1
                  courseNum = 1
            else:
              if blockIndex == 8:
                # Class does not exists
                # Create new class in first available slot 
                blockNum = 1
                while True:
                  newBlock = f"block{blockNum}"
                  if cname not in running[newBlock] and len(running[newBlock]) < blockClassLimit:
                    if student["schedule"][newBlock] == "": # Add student to class
                      running[newBlock][cname] = {
                        "name": courses[currentCourse]["name"],
                        "students": [student["name"]]
                      }
                      student["schedule"][newBlock] = cname
                      break
                    else:
                      if blockNum == 8:
                        # All student classes have been filled
                        break
                      else:
                        blockNum += 1
                        courseNum = 1
                  else:
                    if blockNum == 8:
                      # No room in school for more classes
                      # print(f"\nError 2: No more room in school for another class")
                      err2 += 1
                      break
                    else:
                      if (courses[currentCourse]["teachers"] - courseNum) > 0:
                        courseNum += 1
                      else:
                        blockNum += 1
                        courseNum = 1
                break
              else:
                blockIndex += 1
                courseNum = 1
          break
        elif currentCourse not in activeCourses:
          if len(student["requests"]) == 8:
            # This student never had any alternatives
            # How to solve?

            # print(f"\nError 1: No more available classes for student {student['name']}")
            err1 += 1
            generate = False
          else:
            if alternateIndex <= (len(student["requests"]) - 1):
              currentCourse = student["requests"][alternateIndex]
              alternateIndex += 1
            else: # Out of alternatives
              # print(f"\nError 1: No more available classes for student {student['name']}")
              err1 += 1
              generate = False


def generateScheduleV2():
  global err1, err2
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
    index = list(activeCourses)[i]
    classRunCount = math.floor(activeCourses[index]["totalrequests"] / classCap)
    # If there is minReq+ requests left, 1 more class could be run
    if (activeCourses[index]["totalrequests"] % classCap) > minReq: classRunCount += 1
    activeCourses[index]["classRunCount"] = classRunCount

    blockIndex = 1
    while classRunCount > 0:
      block = f"block{blockIndex}"
      if activeCourses[index]["code"] not in running[block] and len(running[block]) < blockClassLimit:
        # Generate class and sub 1 from classRunCount
        running[block][activeCourses[index]["code"]] = {
          "name": activeCourses[index]["name"],
          "students": []
        }
        classRunCount -= 1
      else:
        if blockIndex == 8:
          # No room in school for more classes
          # print(f"\nError 2: No more room in school for another class")
          err2 += 1
          classRunCount = 0
        else:
          blockIndex += 1

  for student in mockStudents:
    alternateOffset = len(student["requests"])-8
    alternateIndex = 8
    for i in range(len(student["requests"])-alternateOffset): # Subtract x classes as they are alternatives
      currentCourse = student["requests"][i]
      generate = True
      while generate:
        # If class is allowed to run
        if currentCourse in activeCourses:
          blockIndex = 1
          getFreeBlock = True
          while getFreeBlock:
            block = f"block{blockIndex}"
            if currentCourse in running[block]:
              if student["schedule"][block] == "": # Add student to class
                if len(running[block][currentCourse]["students"]) < classCap:
                  running[block][currentCourse]["students"].append(student["name"])
                  student["schedule"][block] = courses[currentCourse]["name"]
                  getFreeBlock = False
                else: # Find next available class or create new one
                  if blockIndex == 8:  # No available classes
                    if len(student["requests"]) == 8:
                      # This student never had any alternatives
                      # How to solve?

                      # print(f"\nError 1: No more available classes for student {student['name']}")
                      err1 += 1
                      generate = False
                    else:
                      if alternateIndex <= (len(student["requests"]) - 1):
                        currentCourse = student["requests"][alternateIndex]
                        alternateIndex += 1
                      else: # No more alterantive
                        # print(f"\nError 1: No more available classes for student {student['name']}")
                        err1 += 1
                        generate = False
                    getFreeBlock = False
                  else: blockIndex += 1
              else:
                if blockIndex == 8: # No available classes
                  if len(student["requests"]) == 8:
                    # This student never had any alternatives
                    # How to solve?

                    # print(f"\nError 1: No more available classes for student {student['name']}")
                    err1 += 1
                    generate = False
                  else:
                    if alternateIndex <= (len(student["requests"]) - 1):
                      currentCourse = student["requests"][alternateIndex]
                      alternateIndex += 1
                    else: # No more alternatives
                      # print(f"\nError 1: No more available classes for student {student['name']}")
                      err1 += 1
                      generate = False
                  getFreeBlock = False
                else: blockIndex += 1
            else:
              if blockIndex == 8:
                # Class does not exists
                # Resort to alternative
                if alternateIndex <= (len(student["requests"]) - 1):
                  currentCourse = student["requests"][alternateIndex]
                  alternateIndex += 1
                else: # Out of alternatives
                  # print(f"\nError 1: No more available classes for student {student['name']}")
                  err1 += 1
                  generate = False
                break
              else: blockIndex += 1
          break
        elif currentCourse not in activeCourses:
          if len(student["requests"]) == 8:
            # This student never had any alternatives
            # How to solve?

            # print(f"\nError 1: No more available classes for student {student['name']}")
            err1 += 1
            generate = False
          else:
            if alternateIndex <= (len(student["requests"]) - 1):
              currentCourse = student["requests"][alternateIndex]
              alternateIndex += 1
            else: # Out of alternatives
              # print(f"\nError 1: No more available classes for student {student['name']}")
              err1 += 1
              generate = False
      


if __name__ == '__main__':
  if len(sys.argv) == 1:
    print("Missing argument")
    exit()
  if sys.argv[1].upper() == 'V1':
    if len(sys.argv) == 3:
      try:
        studentsNum = int(sys.argv[2])
      except:
        print("Error parsing number of students")
        exit()
    print("Processing...")
    mockStudents = generateMockStudents(studentsNum)
    generateScheduleV1()
  elif sys.argv[1].upper() == 'V2':
    if len(sys.argv) == 3:
      try:
        studentsNum = int(sys.argv[2])
      except:
        print("Error parsing number of students")
        exit()
    print("Processing...")
    mockStudents = generateMockStudents(studentsNum)
    generateScheduleV2()
  else: 
    print("Invalid argument")
    exit()

  print("\n")
  print(f"Error 1: x{err1}")
  print(f"Error 2: x{err2}")
  print("\n")

  ### Displays the number of students in each class
  # for block in running:
  #   print(f"\nBlock: {block}")
  #   print("=====================")
  #   for cl in running[block]:
  #     name = running[block][cl]["name"]
  #     students = len(running[block][cl]["students"])
  #     print(f"Class: {name} | Students: {students}")

  # Count errors in students schedules
  errors = 0
  for i in range(len(mockStudents)):
    count = 0
    for course in mockStudents[i]["schedule"]:
      if mockStudents[i]["schedule"][course]=="": count+=1
    if count > 0: errors += 1
  
  print(f"{errors}/{studentsNum} student(s) have a issue with their schedule")


  with open("schedule.json", "w") as outfile:
    json.dump(running, outfile, indent=2)

  with open("students.json", "w") as outfile:
    json.dump(mockStudents, outfile, indent=2)

  print("Done")
