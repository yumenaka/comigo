//go:build darwin

#import <Cocoa/Cocoa.h>

extern void comigoWailsTrayShow(void);
extern void comigoWailsTrayQuit(void);

@interface ComigoWailsTrayTarget : NSObject
@property(strong) NSStatusItem *statusItem;
@property(strong) NSMenu *menu;
@end

@implementation ComigoWailsTrayTarget

- (void)statusClicked:(id)sender {
    NSEvent *event = [NSApp currentEvent];
    if (event.type == NSEventTypeRightMouseUp) {
        [self.statusItem popUpStatusItemMenu:self.menu];
        return;
    }
    comigoWailsTrayShow();
}

- (void)openComigo:(id)sender {
    comigoWailsTrayShow();
}

- (void)quitComigo:(id)sender {
    comigoWailsTrayQuit();
}

@end

static ComigoWailsTrayTarget *comigoWailsTrayTarget;

void comigoWailsTrayStart(void *iconBytes, int iconLen, char *tooltip, char *showTitle, char *showTip, char *quitTitle, char *quitTip) {
    NSData *iconData = [NSData dataWithBytes:iconBytes length:iconLen];
    NSString *tooltipString = [NSString stringWithUTF8String:tooltip];
    NSString *showTitleString = [NSString stringWithUTF8String:showTitle];
    NSString *showTipString = [NSString stringWithUTF8String:showTip];
    NSString *quitTitleString = [NSString stringWithUTF8String:quitTitle];
    NSString *quitTipString = [NSString stringWithUTF8String:quitTip];

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
        comigoWailsTrayTarget.statusItem.button.target = comigoWailsTrayTarget;
        comigoWailsTrayTarget.statusItem.button.action = @selector(statusClicked:);
        [comigoWailsTrayTarget.statusItem.button sendActionOn:(NSEventMaskLeftMouseUp | NSEventMaskRightMouseUp)];

        NSMenu *menu = [[NSMenu alloc] init];
        [menu setAutoenablesItems:NO];
        NSMenuItem *showItem = [[NSMenuItem alloc] initWithTitle:showTitleString action:@selector(openComigo:) keyEquivalent:@""];
        showItem.target = comigoWailsTrayTarget;
        showItem.toolTip = showTipString;
        [menu addItem:showItem];
        [menu addItem:[NSMenuItem separatorItem]];
        NSMenuItem *quitItem = [[NSMenuItem alloc] initWithTitle:quitTitleString action:@selector(quitComigo:) keyEquivalent:@""];
        quitItem.target = comigoWailsTrayTarget;
        quitItem.toolTip = quitTipString;
        [menu addItem:quitItem];
        comigoWailsTrayTarget.menu = menu;
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
