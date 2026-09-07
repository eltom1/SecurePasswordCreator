import random 

"""
program 

checkear cantidad de caracteres [x]
checkear que tenga caracteres especiales[x]
checkear cantidad de mayusculas[x]
checkear cantidad de minusculas[x]
checkear cantidad de numeros[x]

ejecucion del programa{
	-primero ingresar la contrasena 
	-analizar cada caracter y guardarlo en una variable para saber que tipo de dato es:
		ya sea: 
				-numero 
				-caracter_especial(@,*,etc))
				-mayuscula
				-minuscula
	si no cumple los requisitos volver a ingresar contrasena sino pasa y la imprime 
}
"""
ABECEDARIO_MAYUSCULAS = ["A","B","C","D","E","F","G","H","I","J","K","L","M","N","O","P","Q","R","S","T","U","V","W","X","Y","Z"]
ABECEDARIO_MINUSCULAS =["a","b","c","d","e","f","g","h","i","j","k","l","m","n","o","p","q","r","s","t","u","v","w","x","y","z"]
NUMEROS = ["0","1","2","3","4","5","6","7","8","9"]
CARACTERES_ESPECIALES =["@","!","$","%","&","#","^","+","-","/","|","\\",">","<"]
BINARIO = ["1","0"]

TODO = ABECEDARIO_MAYUSCULAS+ ABECEDARIO_MINUSCULAS + NUMEROS + CARACTERES_ESPECIALES + BINARIO

def password_create(nro):
	caracteres_elegidos = random.choices(TODO, k=nro)
	pws = "".join(caracteres_elegidos)
	contra = pws
	return contra
#check the passwords exceeds requiered 
def check_pass(psw):
	acum = 0

	#contadores
	numero = 0
	caracter_especial= 0
	mayuscula = 0
	minuscula = 0
	
	nro = len(psw)
	
	es_valida = True

	print(nro * "*")

	for i in range(nro):
		if psw[acum].isdigit():
			numero += 1
		elif psw[acum].isupper():
			mayuscula += 1
		elif psw[acum].islower():
			minuscula += 1
		else: 
			caracter_especial += 1
		acum+= 1

	print(f"\nLa cantidad de numeros es de: {numero}\nLa cantidad de mayus es de: {mayuscula}\nLa cantidad de minusculas es de: {minuscula}\nLa cantidad de caracteres especiales es de: {caracter_especial}\n")

	if nro < 8:
		print("please at least 8 characters in total")
		es_valida = False
	if numero == 0:
		print("please input a number ")
		es_valida = False
	if mayuscula == 0:
		print("please input a upper characters")
		es_valida = False
	if minuscula == 0:
		print("please input a lower characters")
		es_valida = False
	if caracter_especial == 0:
		print("please input a special characters")
		es_valida = False
	

	if es_valida:
		print("your password is perfect\n")
		return True				
	else:
		print("⚠Password rejected. Please try again.\n")
		return False


#function main :)
def main():
	tam = int(input("\ninput your lenght password: "))
	#print("\n")
	password = password_create(tam)	
	print(f"\n{password}\n")
	check_pass(password)

if __name__ == '__main__':
	main()



