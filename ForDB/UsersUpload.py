import psycopg2
from contextlib import closing

insert = '''
UPDATE public.tuser
SET birthday=%s, dep_id=%s, post_id=%s
WHERE fam = %s and "name" = %s and otch = %s;;
'''

conn = psycopg2.connect(dbname='portaldb', user='portaluser', 
                        password='PortalDB2020', host='172.20.0.82')
cursor = conn.cursor()


text = '''SELECT * FROM public.tuser where fam = %s and "name" = %s and otch = %s;'''

post = """ SELECT *
FROM public.tpost; """

post_id = {}
cursor.execute(post)

for row in cursor:
	post_id[row[2]] = row[0]




dep = """ SELECT *
FROM public.tdep; """

dep_id = {}
cursor.execute(dep)

for row in cursor:
	dep_id[row[1]] = row[0]



i = 0
with open('C:\\Users\\Moriska32\\Downloads\\Список для нового портала.csv', 'r') as file:
	for line in file:

		line = line.replace("\n", "")

		items = line.split(";")
		
		
		
		if i not in (0,1):
			try:
				cursor.execute(insert, (items[5],dep_id[items[-2]] ,post_id[items[-1]],items[0],items[1],items[2]))
			except Exception as e:
				print(items, e)
		
				
			
			
			


	

		i += 1
conn.commit()
conn.close()




