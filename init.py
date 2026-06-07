import os
import sys

# Function to find and replace text in files
def replace_text_in_files(modulename, oldname='go-printos-backend-quickstart'):
    for root, dirs, files in os.walk("."):
        if ".git" in root:
            continue
        for file in files:
            if file == "init.sh" or file == "init.py":
                continue            
            filepath = os.path.join(root, file)
            try:
                with open(filepath, 'r', encoding='utf-8') as f:
                    content = f.read()
                new_content = content.replace(oldname, modulename)
                with open(filepath, 'w', encoding='utf-8') as f:
                    f.write(new_content)
                print(f"Updated: {filepath}")
            except Exception as e:
                print(f"Error processing {filepath}: {e}")

# Main program
def main():
    modulename = input("Insert the new name of the module: ").strip()
    print()
    
    yn = input(f"The name of your module will be {modulename}, is this correct? (y/n): ").strip().lower()
    if yn == 'y':
        print("Renaming all references!")
        replace_text_in_files(modulename)
    elif yn == 'n':
        print("Exiting...")
        sys.exit(0)
    else:
        print("Invalid response")
        sys.exit(1)
    
    print("Enabling git hooks")
    os.system('git config core.hooksPath .hooks')

if __name__ == "__main__":
    main()