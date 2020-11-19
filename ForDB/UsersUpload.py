import psycopg2
from contextlib import closing

	

conn = psycopg2.connect(dbname='portaldb', user='portaluser', 
                        password='PortalDB2020', host='172.20.0.82')
cursor = conn.cursor()


text = '''SELECT * FROM public.tuser where fam = %s and "name" = %s and otch = %s;'''
i = 0
with open('C:\\Users\\Moriska32\\Downloads\\Список для нового портала.csv', 'r') as file:
	for line in file:

		items = line.split(";")
		
		if i == 3:
			break
		print(items)
		cursor.execute(text, (items[0],items[1],items[2]))
		for row in cursor:
			print(row)

		i += 1






