# myList = {
#   "code1": {
#     "name": "hello",
#     "value": 3
#   },
#   "code2": {
#     "name": "bye",
#     "value": 1
#   }
# }


# if "code1" in myList:
#   print(True)
#   print(list(myList).index("code2"))

myList = ["a","b","c","d","e","f","g","h","i","j"]
print(len(myList))

for i in range(len(myList)-2): # Subtract last 2 classes as they are alternatives
  print()

# Get first alternative
print(myList[len(myList)-2])