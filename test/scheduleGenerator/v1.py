#!/usr/bin/env python3
from util.courses import mockCourses, activeCourses # Fake courses
from util.mockStudents import generateMockStudents # Generate fake students
import json


'''
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

minReq, median, classCap, blockClassLimit = 18, 24, 30, 12
activeCourses = {}
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

# Currently V1 has an average of 9.4% success rate, with an avverage of 90.6% error rate
def generateScheduleV1(mockStudents, mockCourses):
  # Collect data and calculate schedules
  for student in mockStudents:
    # Tally class request
    for request in student["requests"]:
      mockCourses[request]["totalrequests"] += 1
      mockCourses[request]["studentindexes"].append(mockStudents.index(student))
      # Add course to active list if enough requests
      if mockCourses[request]["totalrequests"] > minReq and mockCourses[request]["code"] not in activeCourses: activeCourses[mockCourses[request]["code"]] = mockCourses[request]

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
            cname = f"{mockCourses[currentCourse]['name']}-{courseNum}"
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
                        generate = False
                    getFreeBlock = False
                  else:
                    if (mockCourses[currentCourse]["teachers"] - courseNum) > 0:
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
                        "name": mockCourses[currentCourse]["name"],
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
                      break
                    else:
                      if (mockCourses[currentCourse]["teachers"] - courseNum) > 0:
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
            generate = False
          else:
            if alternateIndex <= (len(student["requests"]) - 1):
              currentCourse = student["requests"][alternateIndex]
              alternateIndex += 1
            else: # Out of alternatives
              # print(f"\nError 1: No more available classes for student {student['name']}")
              generate = False

  return running


if __name__ == '__main__':
  print("Processing...")

  mockStudents = generateMockStudents(400)
  timetable = {}
  timetable["Version"] = 1
  timetable["timetable"] = generateScheduleV1(mockStudents, mockCourses)

  with open("../output/timetable.json", "w") as outfile:
    json.dump(timetable, outfile, indent=2)

  print("Done")
