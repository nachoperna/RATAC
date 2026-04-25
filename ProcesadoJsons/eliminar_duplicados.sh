# Este script elimina los pdf que tienen su copia en docx y todos los archivos que tienen 2+ versiones quedandose con la ultima
carpeta="./prueba/"
docs=$(find "$carpeta" -name "*.docx")
cant_pdfs=$(find "$carpeta" -name "*.pdf" | wc -l)
cant_copias=0
i=0
while IFS="" read -r pdf; do
      i=$((i+1))
      echo $i
      nombre=$(basename "$pdf" | cut -d '.' -f1)
      x=$(echo "$docs" | grep -c "$nombre")
      if [ $x -gt 0 ]; then
            # echo "$nombre tiene su copia en docx"
            cant_copias=$((cant_copias+1))
            rm "$pdf"
      fi
done < <(find "$carpeta" -name "*.pdf" -type f)

echo "$cant_copias/$cant_pdfs tienen su copia en docx y fueron borrados"
