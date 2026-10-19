import json

with open('cmd/seeder/data_anggota.json', 'r') as f:
    data = json.load(f)

for item in data:
    jabatan = item['jabatan']
    jabatan = jabatan.replace('Manahemen', 'Manajemen')
    jabatan = jabatan.replace('Kepala Departemen', 'Ketua Departemen')
    jabatan = jabatan.replace('Kadiv', 'Ketua Divisi')
    item['jabatan'] = jabatan

with open('cmd/seeder/data_anggota.json', 'w') as f:
    json.dump(data, f, indent=2)

print("Berhasil merapihkan data_anggota.json!")
