mkdir tlsComponents

touch tlsComponents/server.cnf
echo "[ req ]
      default_bits        = 2048
      prompt              = no
      default_md          = sha256
      distinguished_name  = dn
      req_extensions      = v3_ext

      [ dn ]
      C  = EN
      ST = nil
      L  = nil
      O  = MyCompany
      OU = IT
      CN = localhost

      [ v3_ext ]
      subjectAltName = @alt_names

      [ alt_names ]
      DNS.1 = localhost
" > tlsComponents/server.cnf


openssl req -new -newkey rsa:2048 -keyout tlsComponents/ca.key -x509 -sha256 -days 365 -out tlsComponents/ca.crt
openssl genrsa -out tlsComponents/server.key 2048
openssl req -new -key tlsComponents/server.key -out tlsComponents/server.csr -config tlsComponents/server.cnf
openssl req -noout -text -in tlsComponents/server.csr
openssl x509 -req -in tlsComponents/server.csr -CA tlsComponents/ca.crt -CAkey tlsComponents/ca.key \
  -CAcreateserial -out tlsComponents/server.crt -days 365 -sha256 -extfile tlsComponents/server.cnf -extensions v3_ext

