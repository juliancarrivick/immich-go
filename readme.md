# Uploader en masse for Immich

The Immich project fulfills all my requirements for managing my photos:

- Self-hosted
- Open source
- Abundant functionalities
- User experience closely resembling Google Photos
- Machine learning capabilities
- Well-documented API
- Includes an import tool
- Continuously enhanced
- ...

Now, I need to migrate my photos to the new system in bulk. Most of my photos are stored in a NAS directory, while photos taken with my smartphone are in the Google Photos application. Some of them are also stored on the NAS at full resolution.

To completely transition away from the Google Photos service, I must set up an Immich server, import my NAS-stored photos, and merge them with my Google Photos collection. However, there are instances where the same pictures exist in both systems, sometimes with varying resolutions. Of course, I want to keep only the best copy of the photo.

# Introducing `immich-go`

- import from folders.
- import from zipped archives without unzipping them.
- import from Google Photos takeout archives, without unzipping them.
- import only missing files or better files (an delete the inferior copy from the server).
- import photos taken within a date range.
- create albums.
- no installation.

See my [motivation](#motivation) for proposing an alternative to the `immich-cli`.

# Installation
The installation is radically simple. Copy the binary executable matching your system on you PC. Done.
It could run on your PC or on your server.


# Running `immich-go`
The `immich-go` program uses the Immich API. Hence it need the server address and a valid API key.

# Common CLI switches

`-server URL` URL of the Immich service<br>
`-key KEY` A key generated by the user. Uploaded photos will belong to the key's owner.<br>



## The `upload` command
Upload photos and videos from a local folders, or zipped archives


```sh
immich-go -server URL -key KEY upload [-albums] [-GooglePhotos] [-date DATE-RANGE] folder1|zip1 folder2|zip2....  
```

### `upload` switches:
`-albums` albums are created either after folder names or Google Photos albums.<br>
`-GooglePhoto` import from a Google Photos structured archive.<br>
`-date YYYY-MM-DD` select photos taken during this day.<br>
`-date YYYY-MM` select photos taken during this month.<br>
`-date YYYY` select photos taken during this year.<br>
`-date YYYY-MM-DD,YYYY-MM-DD` select photos taken within this date range.<br>


### Example

The following command will import from a Google Photos takeout archive photos taken between the 1st and the 30th of June 2019 and create albums if any.
```sh
immich-go -server=http://mynas:2283 -key=zzV6k65KGLNB9mpGeri9n8Jk1VaNGHSCdoH1dY8jQ 
upload -albums -GooglePhotos -date=2019-06 ~/Download/takeout-20230715T073439Z-001.zip ~/Download/takeout-20230715T073439Z-002.zip             
```


### Merging strategy

The local file is analyzed to get following data:
- file size in bytes
- date of capture took from the takeout metadata, the exif data, or the file name with possible.
The index is made of the file name + the size in the same way used by the immich server.

Digital cameras often generate file names with a sequence of 4 digits, leading to generate duplicated names. If the names matches, the capture date must be compared.

Tests are done in this order
1. the index is found in immich --> the name and the size match. We have the file, don't upload it.
1. the file name is found in immich and...
    1. dates match and immich file is smaller than the file --> Upload it, and discard the inferior file
    1. dates match and immich file is bigger than the file --> We have already a better version. Don't upload the file. 
1. Immich don't have it. --> Upload the file.



# Build form sources

to be done

### Thanks

- a shout-out to the immich team for they stunning project

This program use following 3rd party libraries:
- github.com/gabriel-vasile/mimetype to get exact file type
- github.com/rwcarlsen/goexif to get date of capture from JPEG files
- github.com/ttacon/chalk for having logs nicely colored 



# Motivation

The immich-cli tool does a great for importing a tone of files at full speed. However, I want more. So I write this utility for my onw purpose. Maybe, it could help some

## Where the `immich-CLI` falls short

The built-in CLI offers the following capabilities:
1. Uploading files to the Immich server, with the option to delete the original files.
2. Importing files, which stores a reference of the file on the server while retaining the original files in their original location.
3. Creating albums based on the folder structure.
4. Preventing duplicate creations on the server.

### Advanced Expertise Required

The CLI tool is available within the Immich server container, eliminating the need to install the `Node.js` tool belt on your PC. Editing the `docker-compose.yml` file is necessary to access the host's files and retrieve your photos. Uploading photos from a different PC than the Immich server requires advanced skills.

### Limitations with Google Takeout Data

The Google Photos Takeout service saves your collection as massive zip archives. These archives contain photos and accompanying JSON files with limited information such as date, GPS coordinates, and albums.

After unzipping the archive, you can use the CLI tool to upload its contents. However, certain limitations exist:
- Photos are organized in folders by year and albums.
- Photos may be duplicated across year folders and albums.
- Photos might be compressed, potentially affecting the CLI's duplicate detection when comparing them to previously imported photos with finest details.


# Todo list
- [ ] read metadata
    - [ ] png,mp4
    - [ ] heic
-     


# Wish list
- [ ] binary releases
- [ ] import vs upload flag
- [X] check in the photo doesn't exist on the server before uploading
    - [X] but keep files with the same name: ex IMG_0201.jpg if they aren't duplicates
    - [ ] some files may have different names (ex IMG_00195.jpg and IMAGE_00195 (1).jpg) and are true duplicates
- [X] replace the server photo, if the file to upload is better.
- [ ] delete local file after successful upload
- [ ] upload XMP sidecar files 
- [ ] select or exclude assets to upload by
    - [ ] type photo / video
    - [ ] name pattern
    - [ ] glob expression like ~/photos/\*/sorted/
    - [ ] size
    - [X] date of capture within a date range
- [ ] multithreaded 
- [X] import from local folder
- [ ] create albums based on folder
- [X] import from zip archives without unzipping them
- [X] Import Google takeout
    - [X] handle archives without unzipping them
    - [X] manage multi-zip archives (related files are scattered across all zips)
    - [X] manage duplicates assets (from Photo folder and Albums)
    - [X] don't import trashed files
    - [ ] don't import failed videos
    - [ ] handle Archives, 
    - [ ] handle google albums in immich
    - [ ] option to include photos taken by a partner (the partner may also uses immich for her/his own photos)
- [ ] use tags placed in exif data
- [ ] upload from remote folders
    - [ ] ssh
    - [ ] samba

