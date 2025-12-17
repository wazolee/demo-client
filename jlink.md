Jlink micsoda, használat nélkülözhetővé teszi egy teljes JVM futtatását?

A jlink (pontosabban jlink, kisbetűvel) a Java JDK-ban (Java 9 óta) elérhető parancssori eszköz, amely lehetővé teszi egy egyedi, minimalizált Java runtime image (futásidejű környezet) létrehozását.
Miért hasznos a jlink?
Hagyományosan egy Java alkalmazás futtatásához telepíteni kell egy teljes JRE-t (Java Runtime Environment) vagy JDK-t, ami több száz MB méretű, és tartalmaz rengeteg olyan modult és osztályt, amit az adott alkalmazás soha nem használ (pl. Swing GUI, RMI, CORBA stb.).
A jlink használatával:

Csak azokat a Java platform modulokat (pl. java.base, java.logging stb.) tartalmazza a runtime, amikre az alkalmazásodnak tényleg szüksége van.
Jelentősen kisebb lesz a méret (gyakran 30-100 MB helyett a teljes JRE 200+ MB-jával szemben).
Gyorsabb indulás, kisebb memóriaigény.
Különösen hasznos Docker konténerekben, cloud deploymentben vagy önállóan terjesztett alkalmazásoknál, ahol nem akarsz függeni a célgépen telepített teljes JVM-től.

Válasz a második részre
Igen, pontosan ez a lényege: a jlink használatával nélkülözhetővé teszi egy teljes JVM (JRE/JDK) futtatását, mert létrehoz egy önálló, az alkalmazásodhoz szabott runtime-ot, amit együtt terjeszthetsz az appal. Nem kell külön telepíteni Java-t a célgépen – a custom image tartalmaz mindent, ami kell a futtatáshoz (beleértve a JVM szükséges részeit is).
Hogyan használod?

Először azonosítsd az alkalmazásod függőségeit a jdeps eszközzel (pl. jdeps --list-deps yourapp.jar).
Aztán futtasd a jlinket:
jlink --module-path $JAVA_HOME/jmods --add-modules java.base,java.logging,... --output my-custom-runtime
Az eredmény egy könyvtár (pl. my-custom-runtime), amiben van bin/java, és ezzel indíthatod az appot: ./my-custom-runtime/bin/java -jar yourapp.jar

Ha moduláris az alkalmazásod (module-info.java-val), még egyszerűbb. Sok build tool (Maven, Gradle pluginok) már integrálja ezt, pl. Spring Boot-hoz is van jlink támogatás Docker image-ek optimalizálására.
Ha konkrét példára vagy kódra van szükséged, mondd, és segítek részletesebben!

**van egy cli appom, ami linuxon fut. Fel van véve mint parancs, tehát ha futtatom a shellben: $ myapp akkor fut. ez valójában egy java -jar myapp.jar indulását váltja ki. ezt az appot hogy lehet a jlinkkel megoldani, hogy ne legyen JRE függősége? Vagyis hogy a futtató hostra ne kelljen egy JRE-t vagy JDK-t telepíteni.**

Igen, a jlink tökéletesen alkalmas erre a célra: létrehoz egy önálló, minimalizált Java runtime image-et (beleértve a szükséges JVM komponenseket), amit az alkalmazásoddal együtt terjeszthetsz. Így a futtató Linux hoston nem kell külön JRE-t vagy JDK-t telepíteni – az image tartalmaz mindent, ami kell. A folyamat feltételezi, hogy a fejlesztői gépeden van JDK (Java 9+), és az appod JAR-ja készen áll.
Előfeltételek

Az appod JAR-ja legyen moduláris (azaz legyen module-info.java fájl a JAR-ban, ami definiálja a modult). Ha nem moduláris (pl. régi stílusú JAR), akkor is működik, de extra lépések kellenek a függőségek azonosításához.
Használj JDK-t a buildhez (pl. OpenJDK 21, ami Linuxon könnyen telepíthető).
Ha van külső függőség (pl. third-party JAR-ok), azokat is bele kell venni a module-path-be.

Lépések a jlink használatához

Függőségek azonosítása (jdeps használatával)
Ha az appod nem moduláris, ez kulcsfontosságú: **textjdeps --list-deps myapp.jar**
Ez kilistázza a szükséges Java modulokat (pl. java.base, java.logging stb.). Ha van külső JAR függőség, add hozzá őket:textjdeps --module-path /path/to/deps --list-deps myapp.jar
Custom runtime image létrehozása (jlink)
Futtasd a jlinket a JDK-ból:textjlink --add-modules java.base,java.logging,...  # ide a jdeps-ből kapott modulok, vesszővel elválasztva
      --module-path /path/to/deps  # ha vannak külső JAR-ok/modulok
      --output my-custom-runtime   # output directory neve
      --compress=2                 # opcionális: tömörítés a kisebb méretért
      --no-header-files            # opcionális: felesleges fájlok eltávolítása
      --no-man-pagesHa az appod moduláris, add hozzá a saját modulod nevét az --add-modules-hez (pl. --add-modules myapp.module,java.base,...).
Az eredmény egy my-custom-runtime directory lesz, ami ~30-100 MB (függ az appodtól), és tartalmazza:
bin/java: a minimalizált JVM indító.
A szükséges modulokat és libeket.

Az app integrálása az image-be
Másold a myapp.jar-t a my-custom-runtime-be (pl. egy app subdirectory-be, de nem kötelező).
Teszteld helyben:text./my-custom-runtime/bin/java -jar myapp.jarHa fut, jó vagy.

Terjesztés és parancsként való beállítás Linuxon
Csomagold az egészet egy tar.gz-be:texttar -czf myapp-bundle.tar.gz my-custom-runtime myapp.jar(Ha a JAR bent van a runtime-ban, csak a runtime-ot csomagold.)
A futtató hoston:
Unpackeld: tar -xzf myapp-bundle.tar.gz.
Hozz létre egy wrapper scriptet (pl. myapp nevű fájl) a gyökérben:Bash#!/bin/bash
SCRIPT_DIR=$$ (dirname " $$(realpath "$0")")
"$SCRIPT_DIR/my-custom-runtime/bin/java" -jar "$$ SCRIPT_DIR/myapp.jar" " $$@"Tedd futtathatóvá: chmod +x myapp.
Add hozzá a PATH-hoz: pl. másold /usr/local/bin-be, vagy a user home-ban .bashrc-ba add: export PATH=$PATH:/path/to/myapp-dir.
Most a shellben $ myapp indítja, JRE nélkül.



Tippek és megjegyzések

Ha moduláris az appod: Használhatsz --launcher opciót a jlinkben, hogy automatikusan létrehozzon egy indító scriptet:textjlink --add-modules myapp.module,java.base,...
      --launcher myapp=myapp.module/com.example.MainClass
      --output my-custom-runtimeEkkor a my-custom-runtime/bin/myapp közvetlenül indítható, mint egy natív bináris (de még mindig Java alapú).
Méret optimalizálás: Ha még kisebbet akarsz, nézd meg a jlink pluginokat (pl. --strip-debug), vagy használd a GraalVM-et natív image-re (de az külön tool, nem jlink).
Build tool integráció: Ha Maven/Gradle-t használsz, vannak pluginok (pl. Maven jlink plugin vagy Gradle jlink task), amik automatizálják ezt. Pl. Spring Boot CLI appoknál a spring-boot-jarmode-layertools segít layerekre bontani Dockerhez.
Platformfüggőség: A jlink image platform-specifikus (Linux x64-re buildelj Linuxon). Ha cross-platform kell, külön buildelj Windowsra stb.
Hibaelhárítás: Ha hiányzik modul, adj hozzá az --add-modules-hez. Ha natív libek kellenek (pl. AWT-hez), a jlink automatikusan kezeli, ha a modul bent van.

Ha megadod az appod részleteit (pl. moduláris-e, milyen függőségek), vagy példakódot, finomíthatom a lépéseket!

**ez egy maven project**

Mivel Maven projektben dolgozol, sokkal egyszerűbben és automatikusabban megoldhatod a jlinkkel történő custom runtime létrehozását, mint kézzel parancsokkal. A legegyszerűbb és legajánlottabb megoldás a Badass JLink Plugin használata (hivatalos neve: org.beryx:badass-jlink-plugin), ami kifejezetten Mavenhez készült, és szinte mindent elintéz helyetted.
1. A plugin hozzáadása a pom.xml-hez
Add hozzá ezt a plugin-t a <build><plugins> szekcióba:
XML<plugin>
    <groupId>org.beryx</groupId>
    <artifactId>badass-jlink-plugin</artifactId>
    <version>3.0.0</version>  <!-- ellenőrizd a legfrissebb verziót: https://github.com/beryx/badass-jlink -->
    <extensions>true</extensions>
    <configuration>
        <options>
            <option>--strip-debug</option>
            <option>--compress=2</option>
            <option>--no-header-files</option>
            <option>--no-man-pages</option>
        </options>
        <launcher>
            <name>myapp</name>  <!-- ez lesz a parancs neve -->
            <mainClass>com.example.Main</mainClass>  <!-- a te main osztályod teljes neve -->
            <!-- opcionális: JVM argumentumok -->
            <!-- <jvmArgs> -->
            <!--     <jvmArg>-Xmx512m</jvmArg> -->
            <!-- </jvmArgs> -->
        </launcher>
    </configuration>
</plugin>
2. Ha van külső függőséged (ami szinte mindig van)
A plugin automatikusan felismeri a classpath-on lévő JAR-okat, és bele tudja őket tenni a runtime image-be egy külön „extra” könyvtárba.
Ha mégis manuálisan akarod kezelni, vagy speciális eset van, használhatsz <extraDependencies> szekciót, de általában nincs rá szükség.
3. Build és image létrehozása
Futtasd ezt a parancsot:
Bashmvn clean package jlink:jlink
Az eredmény a következő lesz a target/jlink könyvtárban:

image/ – itt van a teljes custom runtime (bin/java, lib-ek, stb.)
image/bin/myapp – egy közvetlenül futtatható bináris (wrapper script), ami pontosan úgy viselkedik, mint egy natív program

4. Tesztelés helyben
Bash./target/jlink/image/bin/myapp [argumentumok]
Ha ez fut, akkor kész is vagy – semmilyen JRE nem kell a gépen.
5. Terjesztés Linuxra
Két gyakori mód:
A. Egyszerű ZIP/TAR csomag
Bashcd target/jlink/image
zip -r ../../myapp-linux.zip .
A felhasználó kicsomagolja, és közvetlenül futtathatja a bin/myapp-t, vagy berakhatja a PATH-ba.
B. Debian/RPM csomag készítése (opcionális, de profi)
A badass-jlink plugin támogatja a jpackage-ot is, amivel .deb vagy .rpm csomagot készíthetsz:
XML<configuration>
    ...
    <jpackage>
        <type>DEB</type>  <!-- vagy RPM -->
        <name>myapp</name>
        <version>1.0.0</version>
        <installDir>/opt/myapp</installDir>
        <linuxShortcut>true</linuxShortcut>
        <linuxMenuGroup>Utility</linuxMenuGroup>
    </jpackage>
</configuration>
Majd: mvn jlink:jpackage
6. Fontos megjegyzések

Java verzió: A plugin a buildhez használt JDK verziójával készíti az image-et (pl. ha JDK 21-gyel buildelsz, az image is 21-es lesz).
Nem moduláris app esetén is működik: A badass-jlink plugin tökéletesen kezeli a hagyományos (nem moduláris) JAR-okat is.
Main osztály: Győződj meg róla, hogy a <mainClass> helyesen van megadva (ugyanaz, mint a JAR manifestjében).
Friss plugin verzió: 2025 decemberében valószínűleg már van 3.x vagy újabb verzió – mindig ellenőrizd a GitHub oldalon: https://github.com/beryx/badass-jlink

Összefoglalva
A Maven projektben a legegyszerűbb út:

Add hozzá a badass-jlink plugint a pom.xml-be.
Állítsd be a <launcher>-t a saját main osztályoddal és a kívánt parancsnévvel (myapp).
Futtasd: mvn jlink:jlink
A target/jlink/image/bin/myapp közvetlenül futtatható Linuxon, JRE telepítése nélkül.

Ha elküldöd a pom.xml releváns részét vagy a main osztály nevét, pontos konfigurációt írok neked!911msFastrpm csomag kellIgen, Maven projektben könnyen készíthetsz RPM csomagot a jlink-es custom runtime-ból a jpackage tool segítségével (Java 14+ óta elérhető). Ehhez a legjobb egy dedikált Maven plugin: a jpackage-maven-plugin (groupId: com.github.akman, artifactId: jpackage-maven-plugin) vagy alternatívaként a jlink-jpackager-maven-plugin.
Ajánlom a jpackage-maven-plugin-t, mert egyszerű és közvetlenül támogatja az RPM-et.
1. Plugin hozzáadása a pom.xml-hez
Add hozzá a <build><plugins> szekcióba:
XML<plugin>
    <groupId>com.github.akman</groupId>
    <artifactId>jpackage-maven-plugin</artifactId>
    <version>0.1.6</version>  <!-- ellenőrizd a legfrissebbet: https://github.com/akman/jpackage-maven-plugin -->
    <executions>
        <execution>
            <id>jpackage</id>
            <phase>package</phase>
            <goals>
                <goal>jpackage</goal>
            </goals>
            <configuration>
                <!-- Alapvető beállítások -->
                <name>myapp</name>  <!-- a csomag neve -->
                <version>${project.version}</version>  <!-- pl. 1.0.0 -->
                <mainClass>com.example.Main</mainClass>  <!-- a te main osztályod -->

                <!-- Linux RPM specifikus -->
                <type>rpm</type>  <!-- ez készít RPM-et -->
                <installDir>/opt/myapp</installDir>  <!-- telepítési könyvtár, pl. /opt/myapp -->
                <linuxShortcut>true</linuxShortcut>  <!-- desktop shortcut -->
                <linuxMenuGroup>Utility</linuxMenuGroup>  <!-- menü kategória -->

                <!-- Metadata az RPM-hez -->
                <description>A CLI alkalmazásom leírása</description>
                <vendor>A Te Céged</vendor>
                <copyright>Copyright 2025</copyright>
                <linuxRpmLicenseType>GPL</linuxRpmLicenseType>  <!-- licenc típus, pl. GPL, Proprietary -->
                <linuxPackageRelease>1</linuxPackageRelease>  <!-- release szám -->

                <!-- Opcionális: ikon, licenc fájl stb. -->
                <!-- <icon>src/main/resources/myapp.ico</icon> -->  <!-- Linuxon .png ajánlott -->
                <!-- <licenseFile>LICENSE.txt</licenseFile> -->

                <!-- jlink optimalizálás (automatikusan futtat jlinket, ha kell) -->
                <jlinkOptions>
                    <option>--strip-debug</option>
                    <option>--compress=2</option>
                    <option>--no-header-files</option>
                    <option>--no-man-pages</option>
                </jlinkOptions>
            </configuration>
        </execution>
    </executions>
</plugin>
2. Előfeltételek a build gépen (Linuxon)

JDK 14+ (ajánlott 21+), mert a jpackage benne van.
RPM build toolok: telepítsd rpm-build-et (pl. sudo dnf install rpm-build Fedora/RHEL-en, vagy sudo apt install rpm Debian-on, de RPM készítéshez rpm-build kell).
A buildet Linuxon kell futtatni az RPM-hez (platform-specifikus).

3. Build parancs
Bashmvn clean package
Az eredmény a target/ könyvtárban lesz: myapp-1.0.0-1.x86_64.rpm (vagy hasonló név).
4. Tesztelés és telepítés

Telepítés: sudo rpm -i target/myapp-*.rpm
Futtatás: a launcher általában /opt/myapp/bin/myapp vagy shortcut a menüben.
Eltávolítás: sudo rpm -e myapp

Alternatíva: jlink-jpackager-maven-plugin
Ha több kontroll kell (pl. nem-moduláris JAR-okhoz jobban kezeli):

groupId: de.ecsec, artifactId: jlink-jpackager-maven-plugin
Hasonló konfiguráció, `<type>rpm</type>-pel.

Tippek

Ha ikon kell: használj .png-t Linuxra, és add hozzá `<icon>path/to/icon.png</icon>.
Ha service-ként fusson (pl. systemd): manuálisan hozzáadhatsz post-install scriptet a plugin <installerOptions>-szel, vagy külön RPM spec-et használj.
Méret: az RPM ~30-100 MB lesz, attól függően, milyen modulokat vesz be a jlink.
Ha hiba van (pl. hiányzó függőség): nézd a console outputot, gyakran rpm-build hiányzik.

Ha elküldöd a main osztály nevét, verziót vagy pom.xml részletet, pontosabb configot adok! Ha másik plugin kell, mondd.
