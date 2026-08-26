//go:build darwin

#import <Cocoa/Cocoa.h>

extern void comigoWailsTrayShow(void);
extern void comigoWailsTrayCopyURL(void);
extern void comigoWailsTrayOpenStore(char *storeURL);
extern void comigoWailsTrayQuit(void);

@interface ComigoWailsTrayTarget : NSObject
@property(strong) NSStatusItem *statusItem;
@property(strong) NSMenu *menu;
@end

@implementation ComigoWailsTrayTarget

- (void)openComigo:(id)sender {
    comigoWailsTrayShow();
}

- (void)copyReaderURL:(id)sender {
    comigoWailsTrayCopyURL();
}

- (void)openStore:(NSMenuItem *)sender {
    comigoWailsTrayOpenStore((char *)[sender.representedObject UTF8String]);
}

- (void)quitComigo:(id)sender {
    comigoWailsTrayQuit();
}

@end

static ComigoWailsTrayTarget *comigoWailsTrayTarget;

void comigoWailsTrayStart(void *iconBytes, int iconLen, char *tooltip,
        char *showTitle, char *showTip, char *copyTitle, char *copyTip,
        char *openDirTitle, char *openDirTip, char **storeURLs, int storeCount,
        char *extraTitle, char *extraTip, char *quitTitle, char *quitTip, char *versionTitle) {
    NSData *iconData = [NSData dataWithBytes:iconBytes length:iconLen];
    NSString *tooltipString = [NSString stringWithUTF8String:tooltip];
    NSString *showTitleString = [NSString stringWithUTF8String:showTitle];
    NSString *showTipString = [NSString stringWithUTF8String:showTip];
    NSString *copyTitleString = [NSString stringWithUTF8String:copyTitle];
    NSString *copyTipString = [NSString stringWithUTF8String:copyTip];
    NSString *openDirTitleString = [NSString stringWithUTF8String:openDirTitle];
    NSString *openDirTipString = [NSString stringWithUTF8String:openDirTip];
    NSString *extraTitleString = [NSString stringWithUTF8String:extraTitle];
    NSString *extraTipString = [NSString stringWithUTF8String:extraTip];
    NSString *quitTitleString = [NSString stringWithUTF8String:quitTitle];
    NSString *quitTipString = [NSString stringWithUTF8String:quitTip];
    NSString *versionTitleString = [NSString stringWithUTF8String:versionTitle];
    NSMutableArray<NSString *> *storeURLStrings = [NSMutableArray arrayWithCapacity:storeCount];
    for (int i = 0; i < storeCount; i++) {
        [storeURLStrings addObject:[NSString stringWithUTF8String:storeURLs[i]]];
    }

    dispatch_async(dispatch_get_main_queue(), ^{
        if (comigoWailsTrayTarget == nil) {
            comigoWailsTrayTarget = [[ComigoWailsTrayTarget alloc] init];
        }
        if (comigoWailsTrayTarget.statusItem == nil) {
            comigoWailsTrayTarget.statusItem = [[NSStatusBar systemStatusBar] statusItemWithLength:NSVariableStatusItemLength];
        }

        NSImage *image = [[NSImage alloc] initWithData:iconData];
        [image setSize:NSMakeSize(16, 16)];
        image.template = NO;
        comigoWailsTrayTarget.statusItem.button.image = image;
        comigoWailsTrayTarget.statusItem.button.toolTip = tooltipString;

        NSMenu *menu = [[NSMenu alloc] init];
        [menu setAutoenablesItems:NO];
        NSMenuItem *showItem = [[NSMenuItem alloc] initWithTitle:showTitleString action:@selector(openComigo:) keyEquivalent:@""];
        showItem.target = comigoWailsTrayTarget;
        showItem.toolTip = showTipString;
        [menu addItem:showItem];

        NSMenuItem *copyItem = [[NSMenuItem alloc] initWithTitle:copyTitleString action:@selector(copyReaderURL:) keyEquivalent:@""];
        copyItem.target = comigoWailsTrayTarget;
        copyItem.toolTip = copyTipString;
        [menu addItem:copyItem];

        NSMenuItem *openDirItem = [[NSMenuItem alloc] initWithTitle:openDirTitleString action:nil keyEquivalent:@""];
        openDirItem.toolTip = openDirTipString;
        openDirItem.enabled = storeURLStrings.count > 0;
        NSMenu *openDirMenu = [[NSMenu alloc] init];
        for (NSString *storeURL in storeURLStrings) {
            NSMenuItem *storeItem = [[NSMenuItem alloc] initWithTitle:storeURL action:@selector(openStore:) keyEquivalent:@""];
            storeItem.target = comigoWailsTrayTarget;
            storeItem.representedObject = storeURL;
            [openDirMenu addItem:storeItem];
        }
        openDirItem.submenu = openDirMenu;
        [menu addItem:openDirItem];

        [menu addItem:[NSMenuItem separatorItem]];
        // “其他”保持倒数第二，版本作为禁用子项；退出始终在最底部。
        NSMenuItem *extraItem = [[NSMenuItem alloc] initWithTitle:extraTitleString action:nil keyEquivalent:@""];
        extraItem.toolTip = extraTipString;
        NSMenu *extraMenu = [[NSMenu alloc] init];
        NSMenuItem *versionItem = [[NSMenuItem alloc] initWithTitle:versionTitleString action:nil keyEquivalent:@""];
        versionItem.enabled = NO;
        [extraMenu addItem:versionItem];
        extraItem.submenu = extraMenu;
        [menu addItem:extraItem];

        NSMenuItem *quitItem = [[NSMenuItem alloc] initWithTitle:quitTitleString action:@selector(quitComigo:) keyEquivalent:@""];
        quitItem.target = comigoWailsTrayTarget;
        quitItem.toolTip = quitTipString;
        [menu addItem:quitItem];
        comigoWailsTrayTarget.menu = menu;
        // 使用原生菜单属性，让左右键打开同一菜单；恢复窗口是首项。
        comigoWailsTrayTarget.statusItem.menu = menu;
    });
}

void comigoWailsTrayStop(void) {
    dispatch_async(dispatch_get_main_queue(), ^{
        if (comigoWailsTrayTarget.statusItem != nil) {
            [[NSStatusBar systemStatusBar] removeStatusItem:comigoWailsTrayTarget.statusItem];
            comigoWailsTrayTarget.statusItem = nil;
        }
        comigoWailsTrayTarget = nil;
        [NSApp setActivationPolicy:NSApplicationActivationPolicyRegular];
    });
}

void comigoWailsTraySetWindowVisible(int visible) {
    dispatch_async(dispatch_get_main_queue(), ^{
        if (visible) {
            [NSApp setActivationPolicy:NSApplicationActivationPolicyRegular];
            [NSApp activateIgnoringOtherApps:YES];
            return;
        }
        [NSApp setActivationPolicy:NSApplicationActivationPolicyAccessory];
        dispatch_after(dispatch_time(DISPATCH_TIME_NOW, (int64_t)(0.1 * NSEC_PER_SEC)), dispatch_get_main_queue(), ^{
            [NSApp setActivationPolicy:NSApplicationActivationPolicyAccessory];
        });
    });
}
