import os
import subprocess

def run_tool(command):
    """تشغيل أداة والحصول على النتيجة"""
    try:
        result = subprocess.check_output(command, shell=True, text=True)
        return result.splitlines()
    except subprocess.CalledProcessError as e:
        print(f"Error running command: {command}")
        return []

def remove_duplicates(subdomains):
    """حذف التكرارات"""
    return sorted(set(subdomains))

def run_nuclei(targets_file):
    """تشغيل Nuclei لفحص الثغرات"""
    print("[*] Running Nuclei...")
    nuclei_cmd = f"nuclei -l {targets_file} -o nuclei_results.txt"
    subprocess.call(nuclei_cmd, shell=True)
    print("[*] Nuclei scan completed! Results saved to nuclei_results.txt")

def run_masscan(targets_file):
    """تشغيل Masscan لفحص المنافذ"""
    print("[*] Running Masscan...")
    masscan_cmd = f"masscan -iL {targets_file} -p1-65535 --rate=1000 -oG masscan_results.txt"
    subprocess.call(masscan_cmd, shell=True)
    print("[*] Masscan scan completed! Results saved to masscan_results.txt")

def run_openvas(targets_file):
    """تشغيل OpenVAS لفحص الثغرات العميق"""
    print("[*] Running OpenVAS...")
    # فرض أنك تقوم بتشغيل OpenVAS عبر Docker
    openvas_cmd = f"echo 'OpenVAS placeholder for: {targets_file}'"
    subprocess.call(openvas_cmd, shell=True)
    print("[*] OpenVAS scan completed!")

def main(domain):
    print(f"Running reconnaissance for: {domain}")

    # قائمة لحفظ النتائج
    subdomains = []

    # 1. تشغيل Subfinder
    print("[*] Running Subfinder...")
    subfinder_cmd = f"subfinder -d {domain} -silent"
    subdomains += run_tool(subfinder_cmd)

    # 2. تشغيل Amass
    print("[*] Running Amass...")
    amass_cmd = f"amass enum -passive -d {domain}"
    subdomains += run_tool(amass_cmd)

    # 3. تشغيل Assetfinder
    print("[*] Running Assetfinder...")
    assetfinder_cmd = f"assetfinder --subs-only {domain}"
    subdomains += run_tool(assetfinder_cmd)

    # حذف التكرارات
    print("[*] Removing duplicates...")
    unique_subdomains = remove_duplicates(subdomains)

    # حفظ النتائج إلى ملف
    output_file = f"{domain}_subdomains.txt"
    with open(output_file, "w") as f:
        f.write("\n".join(unique_subdomains))
    print(f"[*] Subdomain enumeration completed! Results saved to {output_file}")

    # 4. فحص الثغرات باستخدام Nuclei
    run_nuclei(output_file)

    # 5. فحص المنافذ باستخدام Masscan
    run_masscan(output_file)

    # 6. فحص الثغرات العميقة باستخدام OpenVAS
    run_openvas(output_file)

if __name__ == "_main_":
    target_domain = input("Enter the domain: ")
    main(target_domain)