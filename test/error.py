import json

f = open('./output/students.json')
students = json.load(f)
f.close()

couldnt_resolve = 0
missing_classes = 0
acceptable_missing_classes = 0

for student in students:
  blocks = [student["schedule"][block] for block in student["schedule"]]
  conflicts = sum(1 for b in blocks if len(b)>1)
  if conflicts == 0 and student["classes"] == student["expectedClasses"]: continue
  if conflicts > 0: couldnt_resolve += 1
  if (student["expectedClasses"] - 2) <= student["classes"] < student["expectedClasses"]: acceptable_missing_classes += 1
  if student["classes"] < (student["expectedClasses"] - 2): missing_classes += 1

print(f"Couldn't resolve            :  {couldnt_resolve}/{len(students)} - {round((couldnt_resolve/len(students))*100, 2)}%")
print(f"Missing classes             :  {missing_classes}/{len(students)}  - {round((missing_classes/len(students))*100, 2)}%")
print(f"Acceptable Missing classes  :  {acceptable_missing_classes}/{len(students)} - {round((acceptable_missing_classes/len(students))*100, 2)}%")
