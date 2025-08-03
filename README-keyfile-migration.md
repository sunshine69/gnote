Gnote support database encryption. In the same directory wehre the gnote sqlite database file is there will be afile with the same name and with .key suffix.

If you run gnote without any argument, the db file will be your user home dir / .gnote.db and your key file will be home dir / .gnote.db.key.

The key file is an base64 of encrypted data. The data is used to encrypt teh database which is a 32 bytes random and converted to hex whihc will be 64 chars. That 64 chars is encrypted using your passphrase.

When first time start gnote will prompt you to create a passphrase. If no passphrase provided, the db will NOT be encrypted. From now one if you run gnote it will prompt but just hit enter to open.

If you enter a pass phrase, gnote will make a 32 byte random, convert to hex and then use your passphrase to encrypt the data and save it to the key file. Then it uses that 32 byte random for sqlite AES 256 encryption.

The below section is only for handly the key file. It does not impact the security level of sqlite encryption so your raw 32 bytes random data of the key file do not need to change.

The old version using a bit weaker encryption function. Recently I update it thus if you use it then you need to run the migration command otherwise it wont be able to decrypt the key file. On console (in windows run using powershell or cmd.exe is fine). type the path to the gnote binary (or in windows the bundle folder /bin/gnote-windows-amd64.exe). Run like this (replace 'gnote' with actual exe file)

```
gnote cli -command migrate-key-file
```

You can use it as a generic encryption/decryption tool. Run `gnote cli -h` for complete options and help.

Note that there are three encryption config, The old and a bit weaker is 0, and verion 1 is best and used by default.