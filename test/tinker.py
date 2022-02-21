#!/usr/bin/python3
import random
import names
import math
from prettytable import PrettyTable
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
running = []


# Generate n students for mock data
def generateMockStudents(n):
  for _ in range(n):
    newStudent = {
      "name": names.get_full_name(),
      "requests": [], # list of class codes
      "schedule": {}
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
            if currentCourse in running[semester][block]:
              # Add student to class
              if len(running[semester][block][currentCourse]["students"]) == classCap:
                # TODO: Class is full, find next class or generate new class
                break
              else:
                running[semester][block][currentCourse]["students"].append(student["name"])
                student["schedule"][semester][block] = running[semester][block][currentCourse]["name"]
            else:
              if semesterIndex == 2 and blockIndex == 4:
                # TODO: Class does not exist, create class in first available slot
                break
              elif blockIndex == 4: # 4th Block (final block)
                blockIndex = 0
                semesterIndex += 1
              else: blockIndex += 1
          break
        elif currentCourse not in activeCourses:
          if alternateOffset == 0:
            # TODO: No more alternatives, solve problem
            break
          currentCourse = student["requests"][len(student["requests"])-alternateOffset]
          alternateOffset -= 1

    # for course in courses:
    #   blockIndex = 0
    #   semesterIndex = 0
    #   if course["totalrequests"] > minReq:
    #     # Create course query
    #     q = {
    #       "name": course["name"],
    #       "remaining": classCap
    #     }
        
    #     # Generate x amount of courses based on requests number
    #     for j in range(len(course["studentindexes"])):
    #       classRunCount = course["totalrequests"] % 30
    #       if ((course["totalrequests"] / 30) - classRunCount) > minReq:
    #         pass
          
    #       # Add one student to class
    #       q["remaining"] -= 1
          
    #       if q["remaining"] == 0 and (len(course["studentindexes"])-classCap) > minReq:
    #         if course["code"] in running[semesterIndex][blockIndex]:
    #           blockIndex += 1
            
    #         # courseIndex = list(running[semesterIndex][blockIndex]).index(course["code"])
    #         running[semesterIndex][blockIndex][course["code"]] = q
        
    #       mockStudents[i]["schedule"][f"semester{semesterIndex+1}"][f"block{blockIndex+1}"] = course["name"]

        # course["totalRequests"] -= minReq
        # while True:
        #   if q in running["semester1"]["block1"]:
        #     running["semester1"][blockIndex][running["semester1"]["block1"].index(q)][]
        #   else: 
        #     running["semester1"][blockIndex].append(q)

      # for studentIndex in course["studentindexes"]:
      #   if studentIndex == i:
      #     pass

generateMockStudents(400)
generateSchedule()