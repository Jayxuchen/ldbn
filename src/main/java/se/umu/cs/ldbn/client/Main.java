package se.umu.cs.ldbn.client;

import se.umu.cs.ldbn.client.i18n.I18N;
import se.umu.cs.ldbn.client.io.Login;
import se.umu.cs.ldbn.client.io.LoginListener;
import se.umu.cs.ldbn.client.ui.GlassPanel;
import se.umu.cs.ldbn.client.ui.admin.AdministratorWidget;
import se.umu.cs.ldbn.client.ui.ca.CreateAssignmentWidget;
import se.umu.cs.ldbn.client.ui.home.HomeWidget;
import se.umu.cs.ldbn.client.ui.licence.LicenceWidget;
import se.umu.cs.ldbn.client.ui.sa.SolveAssignmentWidget;
import se.umu.cs.ldbn.client.ui.user.UserData;
import se.umu.cs.ldbn.client.ui.window.WindowController;
import se.umu.cs.ldbn.client.utils.Common;

import com.allen_sauer.gwt.dnd.client.PickupDragController;
import com.allen_sauer.gwt.log.client.Log;
import com.google.gwt.core.client.EntryPoint;
import com.google.gwt.core.client.Scheduler;
import com.google.gwt.dom.client.Style.Position;
import com.google.gwt.event.logical.shared.SelectionEvent;
import com.google.gwt.event.logical.shared.SelectionHandler;
import com.google.gwt.user.client.Command;
import com.google.gwt.user.client.DOM;
import com.google.gwt.user.client.Window.Location;
import com.google.gwt.user.client.ui.AbsolutePanel;
import com.google.gwt.user.client.ui.Panel;
import com.google.gwt.user.client.ui.RootPanel;
import com.google.gwt.user.client.ui.TabPanel;

/**
 * Entry point classes define <code>onModuleLoad()</code>.
 */
public final class Main implements EntryPoint,
	SelectionHandler<Integer>, LoginListener {

	public static String VERSION = "1.2.2-151201";

	public static final String MAIN_DIV_ID = "gwtapp";

	public static final String LOADING_DIV_ID = "loading";

	public static final int WIDTH_PX;

	static {
		//set the application width in pixels to fit to the screen
		int screenWidth = Common.getScreenWidth();
		screenWidth -= 60;
		int min = 1024-60;
		int max = 1920-60;
		if (screenWidth < min) {
			screenWidth = min;
		} else if (screenWidth > max) {
			screenWidth = max;
		}

		//see if the user has not defined a different screen width
		String w = Location.getParameter("w");
		if (w != null && !"".equals(w.trim())) {
			try {
				int tmp = Integer.parseInt(w);
				if (tmp >= 640 && tmp <= 4000) {
					screenWidth = tmp;
				}
			} catch (Exception e) {
				System.out.println(e);
			}
		}

		WIDTH_PX = screenWidth;
	}

	public static final String REPLACE_LOADING =
		"<table width=\"" + WIDTH_PX + "px\" border=\"0\" style=\"color: gray; text-decoration:none;\"><tr>" +
		"<td style=\"text-align:left; font-size:x-small\">Copyright &copy; 2008 - 2015 " +
		"<a href=\"mailto:ens07ngv@cs.umu.se\" style=\"color: gray; text-decoration:none;\">" +
		"Nikolay Georgiev</a></td>"+
		"<td style=\"text-align: right; font-size:x-small\">Version: " + VERSION + "</td>"+
		"</tr></table>";


	private static Main instance;

	public static Main get() {
		// instance is created in the onModuleLoad2() method,
		// which should be called only once at the begining
		if (instance == null) {
			throw new IllegalArgumentException("onModuleLoad method must be called first.");
		}
		return instance;
	}

	private PickupDragController dragControll;
	private WindowController windowControll;

	private AbsolutePanel mainPanel;

	private GlassPanel glass;

	private TabPanel tabs;

	private AbsolutePanel tabSA;
	private boolean isTabSALoaded;
	private AbsolutePanel tabCA;
	private boolean isTabCALoaded;
	private boolean isTabAdminLoaded;

	private Main() {
		super();
	}

	/**
	 * This is the entry point method.
	 */
	public void onModuleLoad() {
		Log.setUncaughtExceptionHandler();
		//used by the logger!
		Scheduler.get().scheduleDeferred(new Command() {
			public void execute() {
				onModuleLoad2();
			}
		});
	}

	private void onModuleLoad2() {
	    if (Log.isDebugEnabled()) {
	    	Log.debug("Debug log started.");
	    }
		if (instance == null) {
			instance = this;
		}

		//init I18N
		I18N.get();
		AdministratorWidget.get();
		Login.get().addListener(this);

		mainPanel = new AbsolutePanel();
		mainPanel.setWidth(WIDTH_PX+"px");
		dragControll = new PickupDragController(RootPanel.get(), false);
		dragControll.setBehaviorDragProxy(true);

		windowControll = new WindowController(RootPanel.get());

		//tabs
		tabSA = new AbsolutePanel();
		isTabSALoaded = false;
		tabCA = new AbsolutePanel();
		isTabCALoaded = false;
		isTabAdminLoaded = false;

		tabs = new TabPanel();
		tabs.add(HomeWidget.get(), I18N.constants().homeTab());
		tabs.add(tabSA, I18N.constants().solveTab());
		tabs.add(tabCA, I18N.constants().createTab());
		tabs.add(LicenceWidget.get(), I18N.constants().licenseTab());
		tabs.addSelectionHandler(this);

		tabs.setWidth("100%");
		tabs.selectTab(0);
		mainPanel.add(tabs);
		glass = new GlassPanel(mainPanel);
		DOM.setInnerHTML(RootPanel.get(LOADING_DIV_ID).getElement(), REPLACE_LOADING);
		Panel rp = RootPanel.get(MAIN_DIV_ID);
		/*
		 * add(widget, int, int) method requires the parent elements
		 * also be positioned elements, i.e. they must have a property: position.
		 * This is used by the WindowPanel class in order to add the window in the same
		 * static contex as the mainPanel. Then we give each instance of WindowPanel
		 * z-index: 1, thus ensuring the window is infront of the other content.
		 *
		 * @see http://css-discuss.incutio.com/wiki/Overlapping_And_ZIndex
		 */
		rp.getElement().getStyle().setPosition(Position.RELATIVE);
		rp.add(mainPanel);
	}


	public PickupDragController getDragController() {
		return dragControll;
	}

	public WindowController getWindowController() {
		return windowControll;
	}

	public void showGlassPanelLoading() {
		glass.setLoadingAnimationNextTime(true);
		glass.show();
	}

	public void showGlassPanel() {
		glass.show();
	}

	public void hideGlassPanel() {
		glass.hide();
	}

	public void loadSATab() {
		if(!isTabSALoaded) {
			tabSA.add(SolveAssignmentWidget.get());
			isTabSALoaded = true;
		}
	}

	public void onLoginSuccess() {
		if (UserData.get().isAdmin()) {
			tabs.insert(AdministratorWidget.get(), "Administrators", 3);
			isTabAdminLoaded = true;
		}

	}

	public void onSessionKilled() {
		if (isTabAdminLoaded) {
			tabs.remove(3);
			isTabAdminLoaded = false;
		}
	}

	@Override
	public void onSelection(SelectionEvent<Integer> event) {
		int tabIndex = event.getSelectedItem();
		if(tabIndex == 1) {
			loadSATab();
		} else if (tabIndex == 2 && !isTabCALoaded) {
			tabCA.add(CreateAssignmentWidget.get());
			isTabCALoaded = true;
		}
	}
}
