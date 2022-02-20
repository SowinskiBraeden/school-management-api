#!/usr/bin/python3
import random
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
    "block1": {classCode:{"className":name,"remaining":number},classCode:{"className":name,"remaining":number}},
    "block2": {classCode:{"className":name,"remaining":number},classCode:{"className":name,"remaining":number}},
    "block3": {classCode:{"className":name,"remaining":number},classCode:{"className":name,"remaining":number}},
    "block4": {classCode:{"className":name,"remaining":number},classCode:{"className":name,"remaining":number}}
  },
  "semester2": {
    "block1": {classCode:{"className":name,"remaining":number},classCode:{"className":name,"remaining":number}},
    "block2": {classCode:{"className":name,"remaining":number},classCode:{"className":name,"remaining":number}},
    "block3": {classCode:{"className":name,"remaining":number},classCode:{"className":name,"remaining":number}},
    "block4": {classCode:{"className":name,"remaining":number},classCode:{"className":name,"remaining":number}}
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
      "name": "someName",
      "requests": [],
      "schedule": {}
    }
    # Get list of random class choices with no repeats
    # 8 primary choices, 2 secondary choices
    courseSelection = random.sample(range(0, len(courses), 10))
    for courseNum in courseSelection:
      newStudent["requests"] = courses[courseNum]["code"]
    mockStudents.append(newStudent)


def generateSchedule():
  # Collect data and calculate schedules
  for student in mockStudents:
    # Tally class request
    for request in student["requests"]:
      courses[request]["totalrequests"] += 1
      courses[request]["studentindexes"].append(i)
      # Add course to active list if enough requests
      if courses[request]["totalRequests"] > minReq and courses[request] not in activeCourses: activeCourses.append(courses[request])

  for student in mockStudents:
    blockIndex = 0
    semesterIndex = 0
    for i in range(len(student["requests"])-2): # Subtract last 2 classes as they are alternatives
      currentClass = student
      if student["requests"][i] in activeCourses:
        pass
      elif student["requests"][i] not in activeCourses:
        # Get first alternative
        student["requests"][len(student["requests"])-2]

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

      # for studentIndex in course["studentindexs"]:
      #   if studentIndex == i:
      #     pass

generateMockStudents(400)