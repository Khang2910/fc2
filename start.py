fnames = []

for fname in fnames:
    with open(fname) as f:
        content = f.read()
    with open(fname, 'w') as f:
        new_content = content.replace('192.168.2.135', 'newip')
        f.write(new_content)

  
