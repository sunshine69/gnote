#!/usr/bin/python

# shell in windows yuk, when reading the dll_files.lst and convert into the cp
# command it just add some $ char and thus file not found.
# python shutil works and I dont bother hacking bash to work

#export PATH=/usr/local/bin:/usr/bin:/bin:/opt/bin:/c/Windows/System32:/c/Windows:/c/Windows/System32/Wbem:/c/Windows/System32/WindowsPowerShell/v1.0/:/usr/bin/site_perl:/usr/bin/vendor_perl:/usr/bin/core_perl

import shutil, os, sys, re
from glob import glob

os.system("""
export PATH=/usr/local/bin:/usr/bin:/bin:/opt/bin:/c/Windows/System32:/c/Windows:/c/Windows/System32/Wbem:/c/Windows/System32/WindowsPowerShell/v1.0/:/usr/bin/site_perl:/usr/bin/vendor_perl:/usr/bin/core_perl

cd /c/ansible_install
rm -rf gnote-windows-bundle || true
mkdir gnote-windows-bundle

cd gnote-windows-bundle
mkdir bin lib share

cp -a /mingw64/lib/gdk-pixbuf-2.0 lib/
cp -a /usr/share/glib-2.0 share/
cp -a /mingw64/share/glib-2.0/schemas/* share/glib-2.0/schemas/
cp -a /mingw64/share/gtksourceview-3.0 share/
cp -a /mingw64/share/icons share/
cp /c/ansible_install/gnote/gnote*.exe bin/
""")

data = open("/c/ansible_install/gnote/dll_files.lst", "r").read().splitlines()

# Now copy dlls into the bin dir

for f in data:
    if not f.startswith('#'):
        print(f"copying '{f}'")
        #_f1 = glob(f"/mingw64/bin/{_f}*")[0]
        f1 = f"/mingw64/bin/{f}"
        shutil.copy(f1, "/c/ansible_install/gnote-windows-bundle/bin/")

os.chdir("/c/ansible_install")
archive_name = "gnote-windows-bundle-{VER}".format(VER=os.getenv("BUILD_VERSION", "v0.1"))

if os.path.exists(archive_name + ".zip"):
    os.remove(archive_name + ".zip")

shutil.make_archive(archive_name,
    "zip",
    ".",
    "gnote-windows-bundle"
    )

shutil.rmtree("/c/ansible_install/gnote-windows-bundle")
