def extract_urls(input_file, output_file):
    with open(input_file, 'r', encoding='utf-8') as infile, open(output_file, 'w', encoding='utf-8') as outfile:
        for line in infile:
            # Remove any leading/trailing whitespace characters including newline characters
            stripped_line = line.strip()
            # Check if the line starts with "http://" or "https://"
            if stripped_line.startswith('http://') or stripped_line.startswith('https://'):
                # Write the line to the output file
                outfile.write(stripped_line + '\n')

if __name__ == "__main__":
    input_file = 'xm.txt'    # 输入文件
    output_file = 'xm-url.txt'  # 输出文件
    extract_urls(input_file, output_file)
    print(f'URLs have been extracted from {input_file} to {output_file}.')
