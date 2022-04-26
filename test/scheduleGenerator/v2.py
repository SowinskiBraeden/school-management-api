#!/usr/bin/env python3
from util.courses import mockCourses, activeCourses # Fake courses
from util.mockStudents import generateMockStudents # Generate fake students
import math
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


minReq, median, classCap, blockClassLimit = 18, 24, 30, 12
mockStudents = []
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


# Currently V2 has an average of 0.35% success rate, with an avverage of 99.65% error rate
def generateScheduleV2(mockStudents, mockCourses):
  # Collect data and calculate schedules
  for student in mockStudents:
    # Tally class request
    for request in student["requests"]:
      mockCourses[request]["totalrequests"] += 1
      mockCourses[request]["studentindexes"].append(mockStudents.index(student))
      # Add course to active list if enough requests
      if mockCourses[request]["totalrequests"] > minReq and mockCourses[request]["code"] not in activeCourses: activeCourses[mockCourses[request]["code"]] = mockCourses[request]

  # calculate # of times to run class
  for i in range(len(activeCourses)):
    index = list(activeCourses)[i]
    classRunCount = math.floor(activeCourses[index]["totalrequests"] / classCap)
    # If there is minReq+ requests left, 1 more class could be run
    if (activeCourses[index]["totalrequests"] % classCap) > minReq: classRunCount += 1
    activeCourses[index]["classRunCount"] = classRunCount

    newBlockIndex = 1
    newCourseNum = 1
    while classRunCount > 0:
      block = f"block{newBlockIndex}"
      cname = f"{activeCourses[index]['name']}-{newCourseNum}"
      if cname not in running[block] and len(running[block]) < blockClassLimit:
        # Generate class and sub 1 from classRunCount
        running[block][cname] = {
          "name": activeCourses[index]["name"],
          "students": []
        }
        classRunCount -= 1
      else:
        if newBlockIndex == 8:
          if (activeCourses[index]["teachers"] - newCourseNum) > 0:
            newCourseNum += 1
            newBlockIndex = 1
          else:  
            # No room in school for more classes
            # print(f"\nError 2: No more room in school for another class")
            classRunCount = 0
        else: newBlockIndex += 1

  for student in mockStudents:
    alternateOffset = len(student["requests"])-8
    alternateIndex = 8
    for i in range(len(student["requests"])-alternateOffset): # Subtract x classes as they are alternatives
      currentCourse = student["requests"][i]
      generate = True
      courseNum = 1
      while generate:
        # If class is allowed to run
        if currentCourse in activeCourses:
          blockIndex = 1
          getFreeBlock = True
          cname = f"{mockCourses[currentCourse]['name']}-{courseNum}"
          while getFreeBlock:
            block = f"block{blockIndex}"
            if cname in running[block]:
              if student["schedule"][block] == "": # Add student to class
                if len(running[block][cname]["students"]) < classCap:
                  running[block][cname]["students"].append(student["name"])
                  student["schedule"][block] = cname
                  getFreeBlock = False
                else: # Find next available class or create new one
                  if blockIndex == 8:  # No available classes
                    if (activeCourses[index]["teachers"] - newCourseNum) > 0:
                      newCourseNum += 1
                      newBlockIndex = 1
                    else:
                      # No room in school for more classes
                      # print(f"\nError 2: No more room in school for another class")
                      classRunCount = 0
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
                  else: blockIndex += 1
              else:
                if blockIndex == 8: # No available classes
                  if len(student["requests"]) == 8:
                    # This student never had any alternatives
                    # How to solve?

                    # print(f"\nError 1: No more available classes for student {student['name']}")
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
                  if (mockCourses[currentCourse]["teachers"] - courseNum) > 0:
                    courseNum += 1
                  else:
                    blockIndex += 1
                    courseNum = 1
            else:
              if blockIndex == 8:
                # Class does not exists
                # Resort to alternative
                if alternateIndex <= (len(student["requests"]) - 1):
                  currentCourse = student["requests"][alternateIndex]
                  alternateIndex += 1
                else: # Out of alternatives
                  # print(f"\nError 1: No more available classes for student {student['name']}")
                  generate = False
                break
              else:
                if (mockCourses[currentCourse]["teachers"] - courseNum) > 0:
                  courseNum += 1
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
  timetable["Version"] = 2
  timetable["timetable"] = generateScheduleV2(mockStudents, mockCourses)
  
  with open("../output/timetable.json", "w") as outfile:
    json.dump(timetable, outfile, indent=2)

  print("Done")
