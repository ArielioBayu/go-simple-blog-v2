ALTER TABLE uploads
    CHANGE COLUMN nama_file original_filename VARCHAR(255) NOT NULL,
    DROP COLUMN nama_file_sistem,
    CHANGE COLUMN type_file mime_type VARCHAR(50) NOT NULL;
