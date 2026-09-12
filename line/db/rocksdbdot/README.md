### rocksdb

#### window

```bash
# in dot/vcpkg.json, config the rocksdb version and features
{
  "name": "go-rocksdb-project",
  "version-string": "1.0.0",
  "dependencies": [
    {
      "name": "rocksdb",
      "features": ["lz4", "snappy", "zstd"]
    }
  ],
  "overrides": [
    {
      "name": "rocksdb",
      "version-string": "10.10.1"
    }
  ],
  "builtin-baseline": "aa40adda5352e87655b8583cfb2451d5e9e276fd"
}

# update vcpkg baseline
vcpkg.exe x-update-baseline --add-initial-baseline
vcpkg.exe install --triplet=x64-mingw-static
# set the environment variables
VCPKG:="../../../vcpkg_installed"
ROCKSDB_INCLUDE="%VCPKG%/x64-mingw-static/include"
ROCKSDB_LIB="%VCPKG%/x64-mingw-static/lib"
```

<!--```bash
scoop install cmake llvm mingw-winlibs # restart cmd or ide to user new gcc
# scoop bucket add dorado https://github.com/chawyehsu/dorado
scoop update
scoop install zlib bzip2 lz4 zstd #snappy
curl -L -O https://github.com/facebook/rocksdb/archive/refs/tags/v10.10.1.zip
tar -xf v10.10.1.zip # unzip v10.10.1.zip
cd rocksdb-10.10.1
mkdir build
cd build
cmake .. -G "MinGW Makefiles" -DCMAKE_C_COMPILER=gcc -DCMAKE_CXX_COMPILER=g++ -DCMAKE_BUILD_TYPE=Release -DROCKSDB_BUILD_SHARED=OFF -DPORTABLE=ON -DWITH_WERROR=OFF -DWITH_GFLAGS=OFF -DWITH_SNAPPY=OFF -DWITH_ZLIB=OFF -DWITH_LZ4=OFF -DWITH_ZSTD=OFF -DCMAKE_CXX_STANDARD=20 -DCMAKE_CXX_STANDARD_REQUIRED=ON -DCMAKE_CXX_FLAGS="-fpermissive -D_GNU_SOURCE -w"
mingw32-make -j$(nproc) rocksdb
# set the environment variables
ROCKSDB_INCLUDE="${HOME}/rocksdb-10.10.1/include"
ROCKSDB_LIB="${HOME}/rocksdb-10.10.1/build"
```-->

### linux

```bash
sudo apt install -y libbz2-dev
sudo apt install -y libbz2-dev liblz4-dev libzstd-dev libsnappy-dev zlib1g-dev
wget https://github.com/facebook/rocksdb/archive/refs/tags/v10.10.1.zip
unzip v10.10.1.zip
cd rocksdb-10.10.1
make static_lib -j$(nproc)
# set the environment variables
export CFLAGS="-I${HOME}/rocksdb-10.10.1/include"
export LDFLAGS="-L${HOME}/rocksdb-10.10.1/lib"

```
